package wok

import (
	"fmt"
	"net/http"
	"testing"
)

func SessionTestRoutes(app *Wok) {
	app.Get("/session/create", createValue)
	app.Get("/session/get", getValue)
}

func createValue(ctx Context) error {
	err := ctx.SetValue("key", "value")
	if err != nil {
		ctx.Resp.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}

func getValue(ctx Context) error {
	val, err := ctx.GetValue("key")
	if err != nil {
		ctx.Resp.WriteHeader(http.StatusInternalServerError)
		return err
	}

	stringVal, ok := val.(string)
	if !ok {
		ctx.Resp.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("value is not a string")
	}
	if stringVal != "value" {
		ctx.Resp.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("value was supposed to be: \"value\", but returned: %s", stringVal)
	}
	return nil
}

func TestSession(t *testing.T) {
	var client http.Client
	t.Run("create value", func(t *testing.T) {
		res, err := client.Get("http://localhost" + addr + "/session/create")
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != 200 {
			t.Fatal(res.Status)
		}
	})

	t.Run("get value", func(t *testing.T) {
		res, err := client.Get("http://localhost" + addr + "/session/get")
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != 200 {
			t.Fatal(res.Status)
		}
	})

}
