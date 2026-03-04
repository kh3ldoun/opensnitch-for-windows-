package api

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/opensnitch/winsnitch/backend/internal/events"
	"github.com/sirupsen/logrus"
)

type Hub struct {
	log   *logrus.Entry
	mu    sync.RWMutex
	conns map[*websocket.Conn]struct{}
	upg   websocket.Upgrader
}

func NewHub(log *logrus.Entry) *Hub {
	return &Hub{
		log:   log,
		conns: map[*websocket.Conn]struct{}{},
		upg:   websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
	}
}

func (h *Hub) Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upg.Upgrade(w, r, nil)
	if err != nil {
		h.log.WithError(err).Warn("ws upgrade failed")
		return
	}
	h.mu.Lock()
	h.conns[conn] = struct{}{}
	h.mu.Unlock()
}

func (h *Hub) Broadcast(evt events.ConnectionEvent) {
	payload, _ := json.Marshal(evt)
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.conns {
		if err := c.WriteMessage(websocket.TextMessage, payload); err != nil {
			h.log.WithError(err).Debug("ws write failed")
			_ = c.Close()
		}
	}
}

func (h *Hub) Close(ctx context.Context) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.conns {
		_ = c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "shutdown"), time.Now().Add(time.Second))
		_ = c.Close()
	}
}
