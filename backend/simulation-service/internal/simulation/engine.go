// Package simulation is required for describing DSR simulation logic
package simulation

import (
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/domain"
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/node"
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/ws"
	"go.uber.org/zap"
)

type Engine struct {
	Nodes map[int]*node.Node
	Hub   *ws.Hub
	log   *zap.Logger
	step  int
}

func NewEngine(hub *ws.Hub, l *zap.Logger) *Engine {
	return &Engine{
		Nodes: make(map[int]*node.Node),
		Hub:   hub,
		log:   l,
		step:  0,
	}
}

func (e *Engine) InitGraph(nodes []node.Node) {
	for _, nd := range nodes {
		nodeActor := node.NewNode(nd.ID, nd.Neighbors, e.Hub, e.log)
		nodeActor.Start()
		e.Nodes[nd.ID] = nodeActor
	}

	e.log.Info("Simulation engine initialized", zap.Int("nodes", len(nodes)))
}

func (e *Engine) StartRouteDiscovery(source, destination int) {
	e.step++
	e.log.Info("Starting route discovery",
		zap.Int("step", e.step),
		zap.Int("source", source),
		zap.Int("dest", destination))

	rreq := &domain.RREQ{
		RequestID:   e.step,
		Source:      source,
		Destination: destination,
		RouteSoFar:  []int{source},
	}

	if node, ok := e.Nodes[source]; ok {
		node.Inbox <- domain.Message{
			Type: domain.MessageRREQ,
			From: source,
			To:   destination,
			RREQ: rreq,
		}
	}
}

func (e *Engine) Stop() {
	for _, n := range e.Nodes {
		n.Stop()
	}
}
