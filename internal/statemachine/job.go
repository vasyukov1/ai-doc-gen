package statemachine

import "github.com/vasyukov1/ai-doc-gen/internal/domain"

var allowed = map[domain.Status]map[domain.Status]bool{
	domain.StatusQueued: {
		domain.StatusRunning:   true,
		domain.StatusFailed:    true,
		domain.StatusCancelled: true,
	},
	domain.StatusRunning: {
		domain.StatusInputReady:     true,
		domain.StatusNeedsInput:     true,
		domain.StatusRetryWait:      true,
		domain.StatusReadyForReview: true,
		domain.StatusCompleted:      true,
		domain.StatusFailed:         true,
		domain.StatusCancelled:      true,
	},
	domain.StatusInputReady: {
		domain.StatusQueued:    true,
		domain.StatusCancelled: true,
	},
	domain.StatusNeedsInput: {
		domain.StatusQueued:    true,
		domain.StatusCancelled: true,
	},
	domain.StatusRetryWait: {
		domain.StatusRunning:   true,
		domain.StatusFailed:    true,
		domain.StatusCancelled: true,
	},
	domain.StatusReadyForReview: {
		domain.StatusRunning:   true,
		domain.StatusCompleted: true,
		domain.StatusCancelled: true,
	},
}

func CanTransition(from, to domain.Status) bool {
	if from == to {
		return true
	}

	return allowed[from][to]
}
