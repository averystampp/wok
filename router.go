package wok

import (
	"encoding/json"
	"net/http"
)

// Default routes. TODO: Create an override or disable method for developers
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

// Context controls ResponseWriter and pointer to Request, used to extend methods
type Context struct {
	Resp http.ResponseWriter
	Req  *http.Request
}

// Syntactic sugar for passing in data and writing a response in JSON
func (c *Context) JSON(data any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	c.Resp.Write(body)
	return nil
}

// Takes a Wok handler and returns a traditional http.HandlerFunc
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
