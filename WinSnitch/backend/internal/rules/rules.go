package rules

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/opensnitch/winsnitch/backend/internal/events"
)

type Rule struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Enabled   bool            `json:"enabled"`
	Action    events.Decision `json:"action"`
	Process   string          `json:"process"`
	Domain    string          `json:"domain"`
	IP        string          `json:"ip"`
	Port      uint16          `json:"port"`
	Protocol  string          `json:"protocol"`
	Temporary bool            `json:"temporary"`
	Priority  int             `json:"priority"`
}

type Store struct {
	path  string
	mu    sync.RWMutex
	rules []Rule
}

func NewStore(programDataDir string) (*Store, error) {
	rulesPath := filepath.Join(programDataDir, "rules.json")
	s := &Store{path: rulesPath, rules: []Rule{}}
	if _, err := os.Stat(rulesPath); err == nil {
		if err := s.load(); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.rules)
}

func (s *Store) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	blob, err := json.MarshalIndent(s.rules, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, blob, 0o644)
}

func (s *Store) All() []Rule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Rule, len(s.rules))
	copy(out, s.rules)
	return out
}

func (s *Store) Upsert(rule Rule) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.rules {
		if s.rules[i].ID == rule.ID {
			s.rules[i] = rule
			return
		}
	}
	s.rules = append(s.rules, rule)
}
