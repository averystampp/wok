package wok

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

func ConextTestRoutes(app *Wok) {
	app.Get("/context/json", jsonhandler)
}

func jsonhandler(ctx Context) error {
	return ctx.JSON(map[string]string{"test": "json"})
}

func TestJSON(t *testing.T) {
	client := http.Client{}
	res, err := client.Get("http://localhost" + addr + "/context/json")
	if err != nil {
		t.Error(err)
		return
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	defer res.Body.Close()

	want := []byte("{\"test\":\"json\"}")
	if !bytes.Equal(b, want) {
		t.Logf("Got JSON response: %s. But wanted %s", string(b), string(want))
		t.FailNow()
	}
}
