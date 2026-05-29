// Package simulation is required for describing DSR simulation logic
package simulation

import (
	"container/list"
	"sync"

	"github.com/ykshvn/reactive-dsr/shared/events"
	"github.com/ykshvn/reactive-dsr/shared/types"
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/domain"
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/node"
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/queue"
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/ws"
	"go.uber.org/zap"
)

type Engine struct {
	Nodes        map[int]*node.NodeActor
	Hub          *ws.Hub
	log          *zap.Logger
	step         int
	messageQueue *list.List
	queueMu      sync.Mutex
}

func NewEngine(hub *ws.Hub, l *zap.Logger) *Engine {
	return &Engine{
		Nodes:        make(map[int]*node.NodeActor),
		Hub:          hub,
		log:          l,
		step:         0,
		messageQueue: list.New(),
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

func (e *Engine) QueueMessage(msg domain.Message) {
	e.queueMu.Lock()
	defer e.queueMu.Unlock()

	e.messageQueue.PushBack(queue.QueuedMessage{
		Msg:       msg,
		Timestamp: e.step,
		Meta:      make(map[string]interface{}),
	})
}

func (e *Engine) Step() (events.Event, error) {
	e.queueMu.Lock()

	if e.messageQueue.Len() == 0 {
		e.queueMu.Unlock()
		return events.Event{Step: -1}, nil
	}

	element := e.messageQueue.Front()
	e.messageQueue.Remove(element)
	queuedMsg := element.Value.(queue.QueuedMessage)
	e.queueMu.Unlock()

	processedEvent := e.processMessage(queuedMsg.Msg)
	e.step++

	return processedEvent, nil
}

func (e *Engine) processMessage(msg domain.Message) events.Event {
	targetNode, exists := e.Nodes[msg.To]
	if !exists {
		e.log.Warn("Target node not found", zap.Int("target", msg.To))
		return events.Event{Step: -1}
	}

	return targetNode.ProcessSingleMessage(msg)
}

func (e *Engine) SendMessage(msg domain.Message) {
	if _, exists := e.Nodes[msg.To]; exists {
		e.QueueMessage(msg)
	} else {
		e.log.Warn("Target node not found", zap.Int("target", msg.To))
	}
}

func (e *Engine) StartStepRouteDiscovery(src, dst int) {
	if src == dst {
		e.log.Warn("Source and destination are the same")
		return
	}

	if src > len(e.Nodes) || dst > len(e.Nodes) {
		e.log.Error("Source or destination are bad")
		return
	}

	e.queueMu.Lock()
	e.messageQueue = list.New()
	e.queueMu.Unlock()

	e.step = 0

	for _, n := range e.Nodes {
		n.ResetState()
	}

	rreq := &domain.RREQ{
		RequestID:   1,
		Source:      src,
		Destination: dst,
		RouteSoFar:  make([]int, 0),
	}

	e.log.Info("Route discovery started",
		zap.Int("source", src),
		zap.Int("dest", dst))

	e.QueueMessage(domain.Message{
		Type: domain.MessageRREQ,
		From: src,
		To:   src,
		RREQ: rreq,
	})
}

func (e *Engine) GetCurrentStep() int { return e.step }
func (e *Engine) GetQueueLength() int {
	e.queueMu.Lock()
	defer e.queueMu.Unlock()
	return e.messageQueue.Len()
}

func (e *Engine) Stop() {
	for _, n := range e.Nodes {
		n.Stop()
	}
}
