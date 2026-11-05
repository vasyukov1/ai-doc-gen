package documents

import (
	"os"
	"strconv"
	"time"

	"github.com/vasyukov1/ai-doc-gen/internal/domain"
)

// DraftBudget sizes the model stage from the profile. A coverage pass makes
// one request per chapter, so larger documents get proportionally more
// admissions and time instead of starving the last chapters.
type DraftBudget struct {
	Chapters int
	Attempts int
	Duration time.Duration
}

const (
	defaultBaseAttempts       = 12
	defaultAttemptsPerChapter = 2
	defaultBaseSeconds        = 240
	defaultSecondsPerChapter  = 45

	// maxContinuationAttempts bounds one automatic pass after the deadline.
	maxContinuationAttempts = 8
)

// ContinuationPassDuration bounds each automatic pass after the deadline.
const ContinuationPassDuration = 90 * time.Second

// MaxContinuationPasses is the number of automatic passes after the deadline.
const MaxContinuationPasses = 2

func DraftBudgetFor(profile, language string) DraftBudget {
	chapters := len(chapterGroups(profile, Plan(profile, language), domain.Evidence{}))
	if chapters < 1 {
		chapters = 1
	}

	return DraftBudget{
		Chapters: chapters,
		Attempts: max(1, envInt("DRAFT_BASE_ATTEMPTS", defaultBaseAttempts)+envInt("DRAFT_ATTEMPTS_PER_CHAPTER", defaultAttemptsPerChapter)*chapters),
		Duration: time.Duration(envInt("DRAFT_BASE_SECONDS", defaultBaseSeconds)+envInt("DRAFT_SECONDS_PER_CHAPTER", defaultSecondsPerChapter)*chapters) * time.Second,
	}
}

func envInt(key string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil || value < 0 {
		return fallback
	}

	return value
}
