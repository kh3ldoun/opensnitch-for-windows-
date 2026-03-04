package multinode

import "sync"

type Node struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Online  bool   `json:"online"`
}

type Manager struct {
	mu    sync.RWMutex
	nodes map[string]Node
}

func NewManager() *Manager { return &Manager{nodes: map[string]Node{}} }

func (m *Manager) Upsert(n Node) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nodes[n.ID] = n
}

func (m *Manager) List() []Node {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Node, 0, len(m.nodes))
	for _, n := range m.nodes {
		out = append(out, n)
	}
	return out
}
