package scraper

// CompanyConfig defines how to scrape a company's career page.
type CompanyConfig struct {
	Name       string
	CareerURL  string
	SearchTerm string // Search term to look for (intern, internship, etc.)
}

// DefaultCompanies returns the list of companies to monitor.
func DefaultCompanies() []CompanyConfig {
	return []CompanyConfig{
		{
			Name:       "Google",
			CareerURL:  "https://www.google.com/about/careers/applications/jobs/results?q=software+intern&location=United+States",
			SearchTerm: "intern",
		},
		{
			Name:       "Amazon",
			CareerURL:  "https://www.amazon.jobs/en/search?base_query=software+intern&loc_query=United+States",
			SearchTerm: "intern",
		},
		{
			Name:       "Uber",
			CareerURL:  "https://www.uber.com/us/en/careers/list/?query=intern%20software&location=USA",
			SearchTerm: "intern",
		},
		{
			Name:       "DoorDash",
			CareerURL:  "https://careers.doordash.com/jobs/search?query=intern",
			SearchTerm: "intern",
		},
	}
}
