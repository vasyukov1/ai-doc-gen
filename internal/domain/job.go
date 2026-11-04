package domain

import "time"

type Status string

const (
	StatusQueued         Status = "queued"
	StatusRunning        Status = "running"
	StatusInputReady     Status = "input_ready"
	StatusNeedsInput     Status = "needs_input"
	StatusRetryWait      Status = "retry_wait"
	StatusReadyForReview Status = "ready_for_review"
	StatusCompleted      Status = "completed"
	StatusFailed         Status = "failed"
	StatusCancelled      Status = "cancelled"
)

type SourceType string

const (
	SourceGitHub SourceType = "github"
	SourceZIP    SourceType = "zip"
)

type Job struct {
	ID                string     `json:"id"`
	Status            Status     `json:"status"`
	Stage             string     `json:"stage"`
	ProfileID         string     `json:"profile_id"`
	Language          string     `json:"language"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	LastHeartbeatAt   *time.Time `json:"last_heartbeat_at,omitempty"`
	NextRetryAt       *time.Time `json:"next_retry_at,omitempty"`
	RecoveryPending   bool       `json:"recovery_pending"`
	ErrorCode         string     `json:"error_code,omitempty"`
	ErrorMessage      string     `json:"error_message,omitempty"`
	RequiredAction    string     `json:"required_action,omitempty"`
	CompletedSections int        `json:"completed_sections"`
	TotalSections     int        `json:"total_sections"`
	LatestMessage     string     `json:"latest_message"`
}

type Input struct {
	SourceType   SourceType
	GitHubURL    string
	UploadPath   string
	OriginalName string
	SHA256       string
	SizeBytes    int64
}

type NewJob struct {
	ID          string
	ProfileID   string
	Language    string
	Description string
	Input       Input
}

type Source struct {
	JobID         string     `json:"job_id"`
	SourceType    SourceType `json:"source_type"`
	RepositoryURL string     `json:"repository_url,omitempty"`
	CommitSHA     string     `json:"commit_sha,omitempty"`
	SourceDigest  string     `json:"source_digest"`
	WorkspacePath string     `json:"-"`
}

type SourceFile struct {
	Path       string `json:"path"`
	SizeBytes  int64  `json:"size_bytes"`
	SHA256     string `json:"sha256,omitempty"`
	ReasonCode string `json:"reason_code"`
	Eligible   bool   `json:"eligible"`
	Selected   bool   `json:"selected"`
}

type Inventory struct {
	Source Source       `json:"source"`
	Files  []SourceFile `json:"files"`
}
