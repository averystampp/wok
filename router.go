package wok

import (
	"net/http"
	"path/filepath"
)

// Wok enforces its own handler to return an error, then wraps it into an
// http handler converter
type Handler func(Context) error

// Takes a Wok handler and returns a traditional http.HandlerFunc
func handlewokfunc(method string, handle Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, ok := pool.Get().(Context)
		if !ok {
			ctx = Context{}
		}

		ctx.reset(w, r)

		if ctx.Req.Method != method {
			ctx.Resp.WriteHeader(http.StatusMethodNotAllowed)
			ctx.SendString(http.StatusText(http.StatusMethodNotAllowed))
			return
		}

		if err := handle(ctx); err != nil {
			ctx.Resp.Write([]byte(err.Error()))
			return
		}

		WokLog.Info(ctx)
		pool.Put(&ctx)
	}

}

func (wok *Wok) Prefix(fix string) {
	wok.prefix = fix
}

// Enforce the POST method for the passed handler
func (wok *Wok) Post(route string, handle Handler) {
	h := handlewokfunc(http.MethodPost, handle)
	wok.mux.Handle(route, h)
}

// Enforce the GET method for the passed handler
func (wok *Wok) Get(route string, handle Handler) {
	h := handlewokfunc(http.MethodGet, handle)
	wok.mux.Handle(wok.prefix+route, h)
}

// Enforce the PATCH method for the passed handler
func (wok *Wok) Patch(route string, handle Handler) {
	h := handlewokfunc(http.MethodPatch, handle)
	wok.mux.Handle(wok.prefix+route, h)
}

// Enforce PUT method for the passed handler
func (wok *Wok) Put(route string, handle Handler) {
	h := handlewokfunc(http.MethodPut, handle)
	wok.mux.Handle(wok.prefix+route, h)
}

// Enforce the OPTIONS method for the passed handler
func (wok *Wok) Options(route string, handle Handler) {
	h := handlewokfunc(http.MethodOptions, handle)
	wok.mux.Handle(wok.prefix+route, h)
}

// Enforce the DELETE method for the passed handler
func (wok *Wok) Delete(route string, handle Handler) {
	h := handlewokfunc(http.MethodDelete, handle)
	wok.mux.Handle(wok.prefix+route, h)
}

func (wok *Wok) Head(route string, handle Handler) {
	h := handlewokfunc(http.MethodHead, handle)
	wok.mux.Handle(wok.prefix+route, h)
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
