package wok

import (
	"sync"
	"testing"
)

var wg sync.WaitGroup

const (
	addr string = ":5000"
)

func TestMain(m *testing.M) {
	wg.Add(1)
	go func() {
		app := NewWok(addr)
		app.WithSession()

		ConextTestRoutes(app)
		SessionTestRoutes(app)

		app.StartWok()
	}()

	m.Run()
	wg.Done()
	wg.Wait()
}
