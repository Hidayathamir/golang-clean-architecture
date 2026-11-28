package websocket

import "sync"

// Sender is the narrow surface usecases depend on to deliver a real-time
// notification, keeping them decoupled from the concrete hub.
type Sender interface {
	Send(userID int64, payload []byte)
}

var _ Sender = &Hub{}

type Hub struct {
	mu      sync.RWMutex
	clients map[int64]map[*Client]struct{}
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[int64]map[*Client]struct{}),
	}
}

func (h *Hub) Register(userID int64, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[userID] == nil {
		h.clients[userID] = make(map[*Client]struct{})
	}
	h.clients[userID][client] = struct{}{}
}

func (h *Hub) Unregister(userID int64, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	set, ok := h.clients[userID]
	if !ok {
		return
	}
	delete(set, client)
	if len(set) == 0 {
		delete(h.clients, userID)
	}
}

// Send delivers payload to every connection of userID. It is best-effort and
// never blocks: a client whose buffer is full is skipped and re-syncs over
// REST. TODO(multi-instance): fan out via Redis pub/sub so connections held by
// other instances receive the message too.
func (h *Hub) Send(userID int64, payload []byte) {
	h.mu.RLock()
	set := h.clients[userID]
	targets := make([]*Client, 0, len(set))
	for client := range set {
		targets = append(targets, client)
	}
	h.mu.RUnlock()

	for _, client := range targets {
		select {
		case client.Send <- payload:
		default:
		}
	}
}

func (h *Hub) ConnectionCount(userID int64) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients[userID])
}
