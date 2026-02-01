package model

import "time"

// Job represents an intern job listing from a company career page.
type Job struct {
	ID           int64     `json:"id"`
	Company      string    `json:"company"`
	Title        string    `json:"title"`
	URL          string    `json:"url"`
	Location     string    `json:"location,omitempty"`
	DiscoveredAt time.Time `json:"discovered_at"`
	Notified     bool      `json:"notified"`
}
