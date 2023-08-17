package wok

import (
	"io"
	"log"
	"net/http"
	"sync"
	"testing"
)

var wg sync.WaitGroup

func index(ctx Context) error {
	data := struct {
		Name string
		Body []string
		Code int
	}{
		Name: "test",
		Body: []string{"this", "is", "the", "body"},
		Code: 200,
	}
	return ctx.JSON(map[string]interface{}{"data": data}, http.StatusOK)
}

func TestServer(t *testing.T) {
	config := Wok{
		Address:  ":3000",
		Database: false,
	}
	app := NewWok(config)

	app.Get("/", index)

	wg.Add(1)
	go func() {
		app.StartWok()
	}()

	res := testing.Benchmark(BenchmarkServer)
	t.Logf("Server made %d requests", res.N)
	t.Logf("took %d nanoseconds per op", res.NsPerOp())
	t.Logf("wrote %d bytes per op", res.AllocedBytesPerOp())
	t.Logf("and %d allocations per op", res.AllocsPerOp())

	wg.Done()
	wg.Wait()
}

func BenchmarkServer(b *testing.B) {
	client := http.Client{}
	for i := 0; i < b.N; i++ {
		res, err := client.Get("http://localhost:3000/")
		if err != nil {
			log.Fatal(err)
			wg.Done()
		}

		if res.StatusCode != http.StatusOK {
			log.Fatal(err)
			wg.Done()
		}

		_, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
	}
}
