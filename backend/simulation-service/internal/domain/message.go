// Package domain is required for domain logic
package domain

type MessageType string

const (
	MessageRREQ MessageType = "RREQ"
	MessageRREP MessageType = "RREP"
)

type RREQ struct {
	RequestID   int   `json:"request_id"`
	Source      int   `json:"source"`
	Destination int   `json:"destination"`
	RouteSoFar  []int `json:"route_so_far"`
}

type RREP struct {
	RequestID   int   `json:"request_id"`
	Source      int   `json:"source"`
	Destination int   `json:"destination"`
	Route       []int `json:"route"`
}

type Message struct {
	Type MessageType `json:"type"`
	From int         `json:"from"`
	To   int         `json:"to"`
	RREQ *RREQ       `json:"rreq,omitempty"`
	RREP *RREP       `json:"rrep,omitempty"`
}
