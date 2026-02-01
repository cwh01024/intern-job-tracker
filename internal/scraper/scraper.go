package scraper

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"intern-job-tracker/internal/model"

	"golang.org/x/net/html"
)

// Scraper fetches and parses job listings from company career pages.
type Scraper struct {
	client *http.Client
}

// NewScraper creates a new scraper with the given HTTP client.
func NewScraper(client *http.Client) *Scraper {
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	return &Scraper{client: client}
}

// ScrapeAll scrapes all default companies.
func (s *Scraper) ScrapeAll() ([]*model.Job, error) {
	return s.ScrapeAllWithConfigs(DefaultCompanies())
}

// ScrapeAllWithConfigs scrapes all given companies.
func (s *Scraper) ScrapeAllWithConfigs(configs []CompanyConfig) ([]*model.Job, error) {
	var allJobs []*model.Job

	for _, config := range configs {
		jobs, err := s.ScrapeCompany(config)
		if err != nil {
			// Log error but continue with other companies
			fmt.Printf("Error scraping %s: %v\n", config.Name, err)
			continue
		}
		allJobs = append(allJobs, jobs...)
	}

	return allJobs, nil
}

// ScrapeCompany scrapes a single company's career page.
func (s *Scraper) ScrapeCompany(config CompanyConfig) ([]*model.Job, error) {
	resp, err := s.client.Get(config.CareerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", config.CareerURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d for %s", resp.StatusCode, config.CareerURL)
	}

	links := parseJobLinks(resp.Body, config.CareerURL, config.SearchTerm)

	var jobs []*model.Job
	for _, link := range links {
		job := &model.Job{
			Company:      config.Name,
			Title:        link.title,
			URL:          link.url,
			DiscoveredAt: time.Now(),
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

type jobLink struct {
	url   string
	title string
}

// parseJobLinks extracts job links from HTML content.
func parseJobLinks(r io.Reader, baseURL string, searchTerm string) []jobLink {
	var links []jobLink
	seen := make(map[string]bool)

	base, _ := url.Parse(baseURL)
	z := html.NewTokenizer(r)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return links
		case html.StartTagToken:
			t := z.Token()
			if t.Data == "a" {
				var href, text string
				for _, attr := range t.Attr {
					if attr.Key == "href" {
						href = attr.Val
					}
				}

				// Get the text content
				if z.Next() == html.TextToken {
					text = strings.TrimSpace(z.Token().Data)
				}

				// Filter for intern positions
				lowerText := strings.ToLower(text)
				if href != "" && text != "" && strings.Contains(lowerText, strings.ToLower(searchTerm)) {
					// Resolve relative URLs
					linkURL, err := url.Parse(href)
					if err != nil {
						continue
					}
					resolvedURL := base.ResolveReference(linkURL).String()

					if !seen[resolvedURL] {
						seen[resolvedURL] = true
						links = append(links, jobLink{url: resolvedURL, title: text})
					}
				}
			}
		}
	}
}
