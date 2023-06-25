package wok

import "net/http"

// includes all the default routes for user creation, login, logout, and return all users
func DefaultRouter(wok *Wok) {
	wok.Post("/user", CreatUserHandle)
	wok.Post("/login", LoginHandle)
	wok.Get("/all", AllUsers)
	wok.Get("/logout", LogoutUser)
	wok.Delete("/delete", DeleteUserHandle)
}

// Wok enforces its own handler to return an error, then wraps it into an
// http handler converter
type Handler func(Context) error

// context is just a struct of the respose writer and request as used by http handlers
type Context struct {
	Resp http.ResponseWriter
	Req  *http.Request
}

type Wok struct {
	address  string
	mux      *http.ServeMux
	tls      bool
	certFile string
	keyFile  string
}

// Return a new Wok server
func NewWok(tls bool, addr, certfile, keyfile string) *Wok {
	return &Wok{
		address:  addr,
		tls:      tls,
		certFile: certfile,
		keyFile:  keyfile,
		mux:      new(http.ServeMux),
	}
}

// Takes a wok handler and returns a traditional http.HandlerFunc
func handlewokfunc(method string, handle Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := Context{
			Resp: w,
			Req:  r,
		}
		if ctx.Req.Method != method {
			ctx.Resp.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
			return
		}

		if err := handle(ctx); err != nil {
			ctx.Resp.Write([]byte(err.Error()))
			return
		}
	}
}

// Enforce the POST method for the passed handler
func (wok *Wok) Post(route string, handle Handler) {
	h := handlewokfunc("POST", handle)
	wok.mux.Handle(route, h)
}

// Enforce the GET method for the passed handler
func (wok *Wok) Get(route string, handle Handler) {
	h := handlewokfunc("GET", handle)
	wok.mux.Handle(route, h)
}

// Enforce the PATCH method for the passed handler
func (wok *Wok) Patch(route string, handle Handler) {
	h := handlewokfunc("PATCH", handle)
	wok.mux.Handle(route, h)
}

// Enforce PUT method for the passed handler
func (wok *Wok) Put(route string, handle Handler) {
	h := handlewokfunc("PUT", handle)
	wok.mux.Handle(route, h)
}

// Enforce the OPTIONS method for the passed handler
func (wok *Wok) Options(route string, handle Handler) {
	h := handlewokfunc("OPTIONS", handle)
	wok.mux.Handle(route, h)
}

// Enforce the DELETE method for the passed handler
func (wok *Wok) Delete(route string, handle Handler) {
	h := handlewokfunc("DELETE", handle)
	wok.mux.Handle(route, h)
}
