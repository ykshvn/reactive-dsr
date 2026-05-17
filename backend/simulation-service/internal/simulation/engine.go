// Package simulation is required for describing DSR simulation logic
package simulation

import (
	"github.com/ykshvn/reactive-dsr/shared/events"
	"github.com/ykshvn/reactive-dsr/shared/types"
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/domain"
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/node"
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/ws"
	"go.uber.org/zap"
)

type Engine struct {
	Nodes map[int]*node.NodeActor
	Hub   *ws.Hub
	log   *zap.Logger
	step  int
}

func NewEngine(hub *ws.Hub, l *zap.Logger) *Engine {
	return &Engine{
		Nodes: make(map[int]*node.NodeActor),
		Hub:   hub,
		log:   l,
		step:  0,
	}
}

func (e *Engine) InitGraph(nodes []types.Node) {
	for _, nd := range nodes {
		nodeActor := node.NewNodeActor(nd.ID, nd.Neighbors, e.Hub, e.SendMessage, e.log)
		nodeActor.Start()
		e.Nodes[nd.ID] = nodeActor
	}

	e.log.Info("Simulation engine initialized", zap.Int("nodes", len(nodes)))
}

func (e *Engine) Step() {
	e.step++
	e.log.Info("Simulation Step", zap.Int("step", e.step))

	e.Hub.BroadcastEvent(events.NewStepEvent(events.EventSimulationStep, e.step, nil))

	for _, actor := range e.Nodes {
		actor.ProcessPendingMessages()
	}
}

func (e *Engine) SendMessage(msg domain.Message) {
	if targetNode, exists := e.Nodes[msg.To]; exists {
		select {
		case targetNode.Inbox <- msg:
		default:
			e.log.Warn("Inbox full", zap.Int("node", msg.To))

		}
	} else {
		e.log.Warn("Target node not found", zap.Int("target", msg.To))
	}
}

func (e *Engine) StartRouteDiscovery(source, destination int) {
	e.step = 0

	e.log.Info("Starting route discovery",
		zap.Int("source", source),
		zap.Int("dest", destination))

	rreq := &domain.RREQ{
		RequestID:   e.step,
		Source:      source,
		Destination: destination,
		RouteSoFar:  []int{source},
	}

	e.SendMessage(domain.Message{
		Type: domain.MessageRREQ,
		From: source,
		To:   destination,
		RREQ: rreq,
	})

	e.Hub.BroadcastEvent(events.NewEvent(events.EventRREQPropagated, events.RREQPayload{
		From:       source,
		To:         destination,
		RouteSoFar: rreq.RouteSoFar,
		RequestID:  rreq.RequestID,
	}))
}

func (e *Engine) GetCurrentStep() int {
	return e.step
}

func (e *Engine) Stop() {
	for _, n := range e.Nodes {
		n.Stop()
	}
}
