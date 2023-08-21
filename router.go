package wok

import (
	"net/http"
	"path/filepath"
)

// Wok enforces its own handler to return an error, then wraps it into an
// http handler converter
type Handler func(Context) error

// Takes a Wok handler and returns a traditional http.HandlerFunc
func (h Handler) handlewokfunc(method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, ok := pool.Get().(Context)
		if !ok {
			ctx = Context{}
		}

		ctx.reset(w, r)

		if ctx.Req.Method != method {
			ctx.Resp.WriteHeader(http.StatusMethodNotAllowed)
			ctx.SendString(http.StatusText(http.StatusMethodNotAllowed))
			ctx.LogWarn(http.StatusText(http.StatusMethodNotAllowed))
			return
		}

		if err := h(ctx); err != nil {
			ctx.SendString(err.Error())
			return
		}
		pool.Put(&ctx)
	}

}

func (wok *Wok) Prefix(fix string) {
	wok.prefix = fix
}

// Enforce the POST method for the passed handler
func (wok *Wok) Post(route string, handle Handler) {
	wok.mux.Handle(wok.prefix+route, handle.handlewokfunc(http.MethodPost))
	wok.mux.Handle(wok.prefix+route+"/", handle.handlewokfunc(http.MethodPost))
}

// Enforce the GET method for the passed handler
func (wok *Wok) Get(route string, handle Handler) {
	wok.mux.Handle(wok.prefix+route, handle.handlewokfunc(http.MethodGet))
	wok.mux.Handle(wok.prefix+route+"/", handle.handlewokfunc(http.MethodGet))
}

// Enforce the PATCH method for the passed handler
func (wok *Wok) Patch(route string, handle Handler) {
	wok.mux.Handle(wok.prefix+route, handle.handlewokfunc(http.MethodPatch))
	wok.mux.Handle(wok.prefix+route+"/", handle.handlewokfunc(http.MethodPatch))
}

// Enforce PUT method for the passed handler
func (wok *Wok) Put(route string, handle Handler) {
	wok.mux.Handle(wok.prefix+route, handle.handlewokfunc(http.MethodPut))
	wok.mux.Handle(wok.prefix+route, handle.handlewokfunc(http.MethodPut))
}

// Enforce the OPTIONS method for the passed handler
func (wok *Wok) Options(route string, handle Handler) {
	wok.mux.Handle(wok.prefix+route, handle.handlewokfunc(http.MethodOptions))
	wok.mux.Handle(wok.prefix+route+"/", handle.handlewokfunc(http.MethodOptions))
}

// Enforce the DELETE method for the passed handler
func (wok *Wok) Delete(route string, handle Handler) {
	wok.mux.Handle(wok.prefix+route, handle.handlewokfunc(http.MethodDelete))
	wok.mux.Handle(wok.prefix+route+"/", handle.handlewokfunc(http.MethodDelete))
}

func (wok *Wok) Head(route string, handle Handler) {
	wok.mux.Handle(wok.prefix+route, handle.handlewokfunc(http.MethodHead))
	wok.mux.Handle(wok.prefix+route+"/", handle.handlewokfunc(http.MethodHead))
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
