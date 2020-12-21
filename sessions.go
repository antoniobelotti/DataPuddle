package main

import "sync"

type Sessions struct {
	lock  sync.RWMutex
	items map[string]string
}

func (s *Sessions) Add(key string, path string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.items == nil {
		s.items = make(map[string]string)
	}

	s.items[key] = path
}

func (s *Sessions) Get(key string) string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.items[key]
}
