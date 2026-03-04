package dnscache

import "sync"

type Cache struct {
	mu sync.RWMutex
	m  map[string]string
}

func New() *Cache {
	return &Cache{m: map[string]string{}}
}

func (c *Cache) Set(ip, domain string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[ip] = domain
}

func (c *Cache) Get(ip string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	d, ok := c.m[ip]
	return d, ok
}
