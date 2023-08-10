package wok

import (
	"fmt"
	"runtime"
)

type Session struct {
	Name   string
	Items  map[string]interface{}
	Secret string
}

func StartSession() *Session {
	return &Session{
		Items: make(map[string]interface{}),
	}
}

// If key already exists it will overwrite the value
func (s *Session) AddItem(key string, val interface{}) {
	s.Items[key] = val
}

func (s *Session) RetrieveItem(key string) {
	val, ok := s.Items[key]
	if ok {
		fmt.Println(val.(string))
	}
}

func (s *Session) DeleteItem(key string) {
	_, ok := s.Items[key]
	if ok {
		delete(s.Items, key)
	}

}

func (s *Session) NewSession() {
	s.Items = nil
	runtime.GC()
}
