// Package graph used for operations on graph
package graph

type Node struct {
	ID        int     `json:"id"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Neighbors []int   `json:"neighbors"`
}

type GraphResponse struct {
	Nodes     []Node `json:"nodes"`
	Edges     []Edge `json:"edges"`
	NodeCount int    `json:"node_count"`
}

type Edge struct {
	From int `json:"from"`
	To   int `json:"to"`
}
