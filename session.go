package wok

import (
	"sync"
)

type Session struct {
	Name  string
	Items map[string]interface{}
	mu    sync.Mutex
}

func StartSession() *Session {
	return &Session{
		Items: make(map[string]interface{}),
	}
}

// If key already exists it will overwrite the value
func (s *Session) AddItem(key string, val interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Items[key] = val
}

func (s *Session) RetrieveItem(key string) interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.Items[key]
	if ok {
		return val
	}
	return nil
}

func (s *Session) DeleteItem(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Items, key)
}

func (s *Session) NewSession() {
	s.mu.Lock()
	defer s.mu.Unlock()
	clear(s.Items)
}
