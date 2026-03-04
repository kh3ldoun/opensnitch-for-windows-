package blocklist

import (
	"bufio"
	"net/http"
	"strings"
	"sync"
)

type Store struct {
	mu      sync.RWMutex
	domains map[string]struct{}
}

func New() *Store { return &Store{domains: map[string]struct{}{}} }

func (s *Store) Update(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	next := map[string]struct{}{}
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			next[strings.ToLower(fields[1])] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	s.domains = next
	s.mu.Unlock()
	return nil
}

func (s *Store) Contains(domain string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.domains[strings.ToLower(domain)]
	return ok
}
