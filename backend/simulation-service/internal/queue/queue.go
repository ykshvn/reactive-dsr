// Package queue is required for step-by-step simulation
package queue

import (
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/domain"
)

type QueuedMessage struct {
	Msg       domain.Message
	Timestamp int
	Meta      map[string]interface{}
}
