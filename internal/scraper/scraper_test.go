package scraper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestScraper_ScrapeCompany(t *testing.T) {
	// Mock server returning HTML with job listings
	mockHTML := `
	<html>
	<body>
		<div class="job">
			<a href="/jobs/123">Software Engineering Intern, Summer 2026</a>
			<span class="location">Mountain View, CA</span>
		</div>
		<div class="job">
			<a href="/jobs/456">Backend Engineering Intern</a>
			<span class="location">New York, NY</span>
		</div>
	</body>
	</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(mockHTML))
	}))
	defer server.Close()

	scraper := NewScraper(&http.Client{})
	config := CompanyConfig{
		Name:       "TestCompany",
		CareerURL:  server.URL,
		SearchTerm: "intern",
	}

	jobs, err := scraper.ScrapeCompany(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(jobs) == 0 {
		t.Error("expected to find jobs")
	}

	// Check that jobs have correct company
	for _, job := range jobs {
		if job.Company != "TestCompany" {
			t.Errorf("expected company TestCompany, got %s", job.Company)
		}
	}
}

func TestScraper_ScrapeAll(t *testing.T) {
	// Mock server
	mockHTML := `<html><body><a href="/job/1">Software Intern</a></body></html>`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(mockHTML))
	}))
	defer server.Close()

	// Create scraper with custom configs pointing to mock server
	scraper := NewScraper(&http.Client{})
	configs := []CompanyConfig{
		{Name: "Company1", CareerURL: server.URL, SearchTerm: "intern"},
		{Name: "Company2", CareerURL: server.URL, SearchTerm: "intern"},
	}

	jobs, err := scraper.ScrapeAllWithConfigs(configs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have jobs from both companies
	companies := make(map[string]bool)
	for _, job := range jobs {
		companies[job.Company] = true
	}

	if len(companies) != 2 {
		t.Errorf("expected jobs from 2 companies, got %d", len(companies))
	}
}

func TestScraper_ParseJobLinks(t *testing.T) {
	html := `
	<a href="/careers/job/123">Software Engineering Intern</a>
	<a href="/careers/job/456">Data Science Intern 2026</a>
	<a href="/about">About Us</a>
	`

	links := parseJobLinks(strings.NewReader(html), "https://example.com", "intern")

	if len(links) != 2 {
		t.Errorf("expected 2 intern links, got %d", len(links))
	}
}

func TestScraper_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	scraper := NewScraper(&http.Client{})
	config := CompanyConfig{
		Name:      "FailingCompany",
		CareerURL: server.URL,
	}

	_, err := scraper.ScrapeCompany(config)
	if err == nil {
		t.Error("expected error for HTTP 500")
	}
}

func TestScraper_Timeout(t *testing.T) {
	// Test with invalid URL (connection refused)
	scraper := NewScraper(&http.Client{})
	config := CompanyConfig{
		Name:      "TimeoutCompany",
		CareerURL: "http://localhost:99999",
	}

	_, err := scraper.ScrapeCompany(config)
	if err == nil {
		t.Error("expected error for connection failure")
	}
}
