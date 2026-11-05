package domain

import "time"

// Connection is encrypted at rest and never returned by the HTTP API.
type Connection struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
	Key      string `json:"key"`
	RPM      int    `json:"rpm"`
	TPM      int    `json:"tpm"`
}

type DraftRun struct {
	Deadline        time.Time
	Attempts        int
	Reason          string
	Disabled        map[string]string
	SectionAttempts map[string]int
	Connections     []Connection
}
