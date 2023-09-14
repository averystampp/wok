package wok

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type Session struct {
	Name  string
	Items map[string]interface{}
	mu    sync.Mutex
}

func newSession() *Session {
	return &Session{
		Items: make(map[string]interface{}),
	}
}

// If key already exists it will overwrite the value
func (s *Session) AddItem(key string, val interface{}) error {
	if s != nil {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.Items[key] = val
		return nil
	}
	return fmt.Errorf("no session")
}

func (s *Session) RetrieveItem(key string) (interface{}, error) {
	if s != nil {
		s.mu.Lock()
		defer s.mu.Unlock()
		val, ok := s.Items[key]
		if ok {
			return val, nil
		}
		return nil, fmt.Errorf("no values associated with the given key")

	}
	return nil, fmt.Errorf("no session")
}

func (s *Session) DeleteItem(key string) error {
	if s != nil {
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.Items, key)
		return nil
	}
	return fmt.Errorf("no session")
}

func (s *Session) SweepSession() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key := range s.Items {
		b, err := base64.RawURLEncoding.DecodeString(key)
		if err != nil {
			return err
		}
		t := &Token{}
		if err := json.Unmarshal(b, t); err != nil {
			return err
		}

		if t.Expires < time.Now().Unix() {
			delete(s.Items, key)
		}
	}
	return nil
}

// NewSession uses clear() a go 1.21 only stdlib function. This is commented out for now because I
// needed to backport to 1.20.7. I currently require other packages that do not support 1.21
// func (s *Session) NewSession() {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	clear(s.Items)
// }
