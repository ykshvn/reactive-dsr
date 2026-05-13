// Package events required for declaring WebSocket events
package events

import "time"

type EventType string

const (
	EventGraphGenerated     EventType = "graph_generated"
	EventRREQPropagated     EventType = "rreq_propagated"
	EventRREPReceived       EventType = "rrep_received"
	EventRouteFound         EventType = "route_found"
	EventSimulationStep     EventType = "simulation_step"
	EventRouteCacheUpdate   EventType = "route_cache_update"
	EventSimulationFinished EventType = "simulation_finished"
)

type Event struct {
	Type      EventType   `json:"type"`
	Step      int         `json:"step,omitempty"`
	Payload   interface{} `json:"payload"`
	Timestamp int64       `json:"timestamp"`
}

type RREQPayload struct {
	From       int   `json:"from"`
	To         int   `json:"to"`
	RouteSoFar []int `json:"route_so_far"`
	RequestID  int   `json:"request_id"`
}

type RREPPayload struct {
	From  int   `json:"from"`
	To    int   `json:"to"`
	Route []int `json:"route"`
}

func NewEvent(eventType EventType, payload interface{}) Event {
	return Event{
		Type:      eventType,
		Timestamp: time.Now().UnixMilli(),
		Payload:   payload,
	}
}

func NewStepEvent(eventType EventType, step int, payload interface{}) Event {
	return Event{
		Type:      eventType,
		Step:      step,
		Payload:   payload,
		Timestamp: time.Now().UnixMilli(),
	}
}
