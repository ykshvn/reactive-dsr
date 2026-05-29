// Package node describes single node logic
package node

import (
	"container/list"
	"fmt"
	"sync"

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
	NodeQueue    *list.List
	queueMu      sync.Mutex
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
		NodeQueue:    list.New(),
		Hub:          hub,
		sendMessage:  sendMessage,
		log:          l,
		stop:         make(chan struct{}),
	}
}

func (a *NodeActor) Start() { go a.run() }
func (a *NodeActor) Stop()  { close(a.stop) }

func (a *NodeActor) ResetState() {
	a.SeenRequests = make(map[string]bool)
	a.RouteCache = make(map[int][]int)
}

func (a *NodeActor) run() {
	for {
		select {
		case <-a.Inbox:

		case <-a.stop:
			return
		}
	}
}

func (a *NodeActor) ProcessSingleMessage(msg domain.Message) events.Event {
	switch msg.Type {
	case domain.MessageRREQ:
		return a.handleRREQStep(msg.RREQ)
	case domain.MessageRREP:
		return a.handleRREPStep(msg.RREP)
	default:
		return events.Event{Step: -1}
	}
}

func (a *NodeActor) handleRREQStep(rreq *domain.RREQ) events.Event {
	requestKey := fmt.Sprintf("%d-%d", rreq.RequestID, rreq.Source)

	rreq.RouteSoFar = append(append([]int(nil), rreq.RouteSoFar...), a.ID)
	if a.SeenRequests[requestKey] {
		return events.NewEvent(events.EventRREQDropped, events.RREQPayload{
			From:       a.ID,
			To:         rreq.Destination,
			RouteSoFar: rreq.RouteSoFar,
			RequestID:  rreq.RequestID,
		})
	}
	a.SeenRequests[requestKey] = true

	var resultEvent events.Event

	if a.ID == rreq.Destination {
		return a.sendRREPStep(rreq)
	}

	for _, neighbor := range a.Neighbors {
		if len(rreq.RouteSoFar) > 1 && neighbor == rreq.RouteSoFar[len(rreq.RouteSoFar)-2] {
			continue
		}

		routeCopy := append([]int(nil), rreq.RouteSoFar...)

		msg := domain.Message{
			Type: domain.MessageRREQ,
			From: a.ID,
			To:   neighbor,
			RREQ: &domain.RREQ{
				RequestID:   rreq.RequestID,
				Source:      rreq.Source,
				Destination: rreq.Destination,
				RouteSoFar:  routeCopy,
			},
		}

		a.sendMessage(msg)
	}

	if resultEvent.Step != -1 {
		resultEvent = events.NewEvent(events.EventRREQProcessed, events.RREQPayload{
			From:       rreq.Source,
			To:         rreq.Destination,
			RouteSoFar: rreq.RouteSoFar,
			RequestID:  rreq.RequestID,
		})
	}

	return resultEvent
}

func (a *NodeActor) sendRREPStep(rreq *domain.RREQ) events.Event {
	rrep := &domain.RREP{
		RequestID:   rreq.RequestID,
		Source:      rreq.Source,
		Destination: rreq.Destination,
		Route:       append([]int(nil), rreq.RouteSoFar...),
	}

	if len(rrep.Route) <= 1 {
		return events.Event{Step: -1}
	}

	nextHop := rrep.Route[len(rrep.Route)-2]

	a.sendMessage(domain.Message{
		Type: domain.MessageRREP,
		From: a.ID,
		To:   nextHop,
		RREP: rrep,
	})

	return events.NewEvent(events.EventRREPGenerated, events.RREPPayload{
		From:  a.ID,
		To:    rreq.Source,
		Route: rrep.Route,
	})
}

func (a *NodeActor) handleRREPStep(rrep *domain.RREP) events.Event {
	if a.ID != rrep.Source && len(rrep.Route) > 1 {
		myself := findIndex(rrep.Route, a.ID)
		if myself > 0 {
			nextHop := rrep.Route[myself-1]

			a.sendMessage(domain.Message{
				Type: domain.MessageRREP,
				From: a.ID,
				To:   nextHop,
				RREP: rrep,
			})

			return events.NewEvent(events.EventRREPForwarded, events.RREPPayload{
				From:  a.ID,
				To:    nextHop,
				Route: rrep.Route,
			})
		}
	}

	if a.ID == rrep.Source {
		return events.NewEvent(events.EventRouteDiscovered, events.RREPPayload{
			From:  rrep.Source,
			To:    rrep.Destination,
			Route: rrep.Route,
		})
	}

	return events.NewEvent(events.EventRREPReceived, events.RREPPayload{
		From:  rrep.Source,
		To:    rrep.Destination,
		Route: rrep.Route,
	})
}

func findIndex(slice []int, target int) int {
	for i, v := range slice {
		if v == target {
			return i
		}
	}
	return -1
}
