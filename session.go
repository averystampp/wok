package wok

import (
	"encoding/base64"
	"encoding/json"
	"sync"
	"time"
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
// needed to backport to 1.20.7. I currently require other packages that are no 1.21 packages
// func (s *Session) NewSession() {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	clear(s.Items)
// }
