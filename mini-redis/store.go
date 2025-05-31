package main

import "sync"

type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

// returning pointer , because its unsafe to copy mutexes
func New() *Store {
	return &Store{

		data: make(map[string]string),
	}
}
func (s *Store) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, ok := s.data[key]
	return value, ok
}
func (s *Store) Set(key, val string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = val

}
func (s *Store) Del(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}
