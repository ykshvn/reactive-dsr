// Package queue is required for step-by-step simulation
package queue

type Queue struct{}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Push() {
}

func (q *Queue) Pop() {
}
