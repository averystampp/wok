package wok

import (
	"io"
	"net/http/httptest"
	"testing"
)

func TestGetHandler(t *testing.T) {
	r := NewRouter()
	r.Get("/", func(ctx Context) error {
		ctx.Resp.Write([]byte("get request"))
		return nil
	})
	server := httptest.NewServer(r)
	defer server.Close()

	resp, err := server.Client().Get(server.URL + "/")
	if err != nil {
		t.Fatal(err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}

func TestPostHandler(t *testing.T) {
	r := NewRouter()
	r.Post("/", func(ctx Context) error {
		ctx.Resp.Write([]byte("post request"))
		return nil
	})
	server := httptest.NewServer(r)
	defer server.Close()

	resp, err := server.Client().Post(server.URL+"/", "text/plain", nil)
	if err != nil {
		t.Fatal(err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}
