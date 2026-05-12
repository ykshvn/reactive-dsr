// Package types required for declaring shared types between services
package types

type Node struct {
	ID        int     `json:"id"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Neighbors []int   `json:"neighbors"`
}

type Edge struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type GraphResponse struct {
	Nodes     []Node `json:"nodes"`
	Edges     []Edge `json:"edges"`
	NodeCount int    `json:"node_count"`
}
