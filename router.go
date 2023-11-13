package wok

import (
	"net/http"
	"path/filepath"
)

// Wok enforces its own handler to return an error, then wraps it into an
// http handler converter
type Handler func(Context) error

type WokRoute struct {
	path    string
	method  string
	wokFunc Handler
}

func (router *WokRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := Context{
		Resp: w,
		Req:  r,
	}

	if r.Method != router.method {
		ctx.Resp.WriteHeader(http.StatusMethodNotAllowed)
		ctx.Resp.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}

	if err := router.wokFunc(ctx); err != nil {
		ctx.Resp.Write([]byte(err.Error()))

		if wokLogger != nil {
			wokLogger.Info(&ctx, "msg")
		}
		return
	}

	if wokLogger != nil {
		wokLogger.Info(&ctx, "msg")
	}
}

// Enforce the POST method for the passed handler
func (wok *Wok) Post(path string, handle Handler) {
	router := &WokRoute{
		wokFunc: handle,
		method:  http.MethodPost,
	}
	wok.mux.Handle(path, router)
}

func (wok *Wok) Get(path string, handle Handler) {
	router := &WokRoute{
		wokFunc: handle,
		method:  http.MethodGet,
	}
	wok.mux.Handle(path, router)
}

func (wok *Wok) Patch(path string, handle Handler) {
	router := &WokRoute{
		wokFunc: handle,
		method:  http.MethodPatch,
	}
	wok.mux.Handle(path, router)
}

func (wok *Wok) Put(path string, handle Handler) {
	router := &WokRoute{
		wokFunc: handle,
		method:  http.MethodPut,
	}
	wok.mux.Handle(path, router)
}

func (wok *Wok) Options(path string, handle Handler) {
	router := &WokRoute{
		wokFunc: handle,
		method:  http.MethodOptions,
	}
	wok.mux.Handle(path, router)
}

func (wok *Wok) Delete(path string, handle Handler) {
	router := &WokRoute{
		wokFunc: handle,
		method:  http.MethodDelete,
	}
	wok.mux.Handle(path, router)
}

func (wok *Wok) Head(path string, handle Handler) {
	router := &WokRoute{
		wokFunc: handle,
		method:  http.MethodHead,
	}
	wok.mux.Handle(path, router)
}

func (wok *Wok) ServeDir(route string, path string) error {
	ab, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	fs := http.FileServer(http.Dir(ab))
	wok.mux.Handle(route+"/", http.StripPrefix(route, fs))
	return nil
}
