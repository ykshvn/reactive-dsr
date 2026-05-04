package graph

import (
	"math"
	"math/rand"
	"slices"
	"time"
)

type Generator struct {
	rand *rand.Rand
}

func NewGenerator() *Generator {
	return &Generator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (g *Generator) GenerateGraph(n int) (*GraphResponse, error) {
	if n < 3 || n > 50 {
		n = 20
	}

	maxDegree := max(n/2, 3)

	nodes := make([]Node, n)
	edges := make([]Edge, 0)

	centerX := 500.0
	centerY := 500.0
	radius := 500.0

	for i := range n {
		angle := 2 * math.Pi * float64(i) / float64(n)
		x := centerX + radius*math.Cos(angle)
		y := centerY + radius*math.Sin(angle)

		nodes[i] = Node{
			ID: i,
			X:  x,
			Y:  y,
			// X:  g.rand.Float64()*800 + 100,
			// Y:  g.rand.Float64()*800 + 100,
		}
	}

	for i := range n {
		next := (i + 1) % n
		nodes[i].Neighbors = append(nodes[i].Neighbors, next)
		nodes[next].Neighbors = append(nodes[next].Neighbors, i)

		edges = append(edges, Edge{From: i, To: next})
	}

	extraEdges := n * 2
	for range extraEdges {
		u := g.rand.Intn(n)
		v := g.rand.Intn(n)

		if u == v || slices.Contains(nodes[u].Neighbors, v) {
			continue
		}
		if len(nodes[u].Neighbors) >= maxDegree || len(nodes[v].Neighbors) >= maxDegree {
			continue
		}

		nodes[u].Neighbors = append(nodes[u].Neighbors, v)
		nodes[v].Neighbors = append(nodes[v].Neighbors, u)
		edges = append(edges, Edge{From: u, To: v})
	}

	for i := range nodes {
		nodes[i].Neighbors = unique(nodes[i].Neighbors)
	}
	return &GraphResponse{
		Nodes:     nodes,
		Edges:     edges,
		NodeCount: n,
	}, nil
}

func unique(slice []int) []int {
	seen := make(map[int]bool)
	result := []int{}
	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}
