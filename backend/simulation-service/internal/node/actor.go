// Package node describes single node logic
package node

import (
	"fmt"

	"github.com/ykshvn/reactive-dsr/shared/events"
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/domain"
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/ws"

	"go.uber.org/zap"
)

type SendMessageFunc func(domain.Message)

type NodeActor struct {
	ID           int
	Neighbors    []int
	RouteCache   map[int][]int
	SeenRequests map[string]bool
	Inbox        chan domain.Message
	Hub          *ws.Hub
	sendMessage  SendMessageFunc
	log          *zap.Logger
	stop         chan struct{}
}

func NewNodeActor(id int, neighbors []int, hub *ws.Hub, sendMessage SendMessageFunc, l *zap.Logger) *NodeActor {
	return &NodeActor{
		ID:           id,
		Neighbors:    neighbors,
		RouteCache:   make(map[int][]int),
		SeenRequests: make(map[string]bool),
		Inbox:        make(chan domain.Message, 50),
		Hub:          hub,
		sendMessage:  sendMessage,
		log:          l,
		stop:         make(chan struct{}),
	}
}

func (a *NodeActor) Start() {
	go a.run()
}

func (a *NodeActor) Stop() {
	close(a.stop)
}

func (a *NodeActor) run() {
	for {
		select {
		case msg := <-a.Inbox:
			a.processMessage(msg)
		case <-a.stop:
			return
		}
	}
}

func (a *NodeActor) ProcessPendingMessages() {
	for {
		select {
		case msg := <-a.Inbox:
			a.processMessage(msg)
		default:
			return
		}
	}
}

func (a *NodeActor) processMessage(msg domain.Message) {
	switch msg.Type {
	case domain.MessageRREQ:
		a.handleRREQ(msg.RREQ)
	case domain.MessageRREP:
		a.handleRREP(msg.RREP)
	}
}

func (a *NodeActor) sendRREP(rreq *domain.RREQ) {
	route := make([]int, len(rreq.RouteSoFar))
	copy(route, rreq.RouteSoFar)

	for i, j := 0, len(route)-1; i < j; i, j = i+1, j-1 {
		route[i], route[j] = route[j], route[i]
	}

	rrep := &domain.RREP{
		RequestID:   rreq.RequestID,
		Source:      rreq.Source,
		Destination: rreq.Destination,
		Route:       route,
	}

	a.log.Info(
		"Route found! Sending RREP",
		zap.Int("node", a.ID),
		zap.Int("to", rreq.Source),
		zap.Any("route", route),
	)

	if len(route) > 1 {
		nextHop := route[1]

		a.sendMessage(
			domain.Message{
				Type: domain.MessageRREP,
				From: a.ID,
				To:   nextHop,
				RREP: rrep,
			},
		)
	}

	a.Hub.BroadcastEvent(events.NewEvent(events.EventRREPReceived, events.RREPPayload{
		From:  a.ID,
		To:    rreq.Source,
		Route: route,
	}))
}

func (a *NodeActor) handleRREQ(rreq *domain.RREQ) {
	requestKey := fmt.Sprintf("%d-%d", rreq.RequestID, rreq.Source)

	if a.SeenRequests[requestKey] {
		return
	}
	a.SeenRequests[requestKey] = true

	if a.ID == rreq.Destination {
		a.sendRREP(rreq)
	}

	newRoute := append([]int{}, rreq.RouteSoFar...)
	newRoute = append(newRoute, a.ID)

	for _, neighbor := range a.Neighbors {
		if neighbor == rreq.RouteSoFar[len(rreq.RouteSoFar)-1] {
			continue
		}

		a.sendMessage(domain.Message{
			Type: domain.MessageRREQ,
			From: a.ID,
			To:   neighbor,
			RREQ: &domain.RREQ{
				RequestID:   rreq.RequestID,
				Source:      rreq.Source,
				Destination: rreq.Destination,
				RouteSoFar:  newRoute,
			},
		})
	}

	a.Hub.BroadcastEvent(events.NewEvent(events.EventRREQPropagated, events.RREQPayload{
		From:       a.ID,
		To:         rreq.Destination,
		RouteSoFar: newRoute,
		RequestID:  rreq.RequestID,
	}))
}

func (a *NodeActor) handleRREP(rrep *domain.RREP) {
	a.log.Info("RREP received",
		zap.Int("node", a.ID),
		zap.Int("destination", rrep.Destination),
		zap.Any("full_route", rrep.Route),
	)

	if a.ID != rrep.Source && len(rrep.Route) > 1 {
		nextHop := rrep.Route[1]

		a.sendMessage(domain.Message{
			Type: domain.MessageRREP,
			From: a.ID,
			To:   nextHop,
			RREP: rrep,
		})

		a.log.Info("Forwarding RREP",
			zap.Int("from", a.ID),
			zap.Int("to", nextHop),
		)
	}

	a.Hub.BroadcastEvent(events.NewEvent(events.EventRREPReceived, events.RREPPayload{
		From:  a.ID,
		To:    rrep.Source,
		Route: rrep.Route,
	}))
}
