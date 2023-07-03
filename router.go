package wok

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
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
	Ctx  context.Context
}

// Syntactic sugar for passing in data and writing a response in JSON
func (ctx *Context) JSON(data any) error {

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	ctx.Resp.Header().Set("Content-Type", "application/json")
	ctx.Resp.Write(body)
	return nil
}

// Set key and value pairs for Ctx
func (ctx *Context) SetKey(key any, val any) {
	ctx.Ctx = context.WithValue(ctx.Ctx, key, val)
}

// Returns value of key as a string
func (ctx *Context) GetKey(key any) string {
	valuefromCtx := ctx.Ctx.Value(key)
	value := fmt.Sprintf("%v", valuefromCtx)
	return value
}

// Takes a Wok handler and returns a traditional http.HandlerFunc
func handlewokfunc(method string, handle Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := Context{
			Resp: w,
			Req:  r,
			Ctx:  context.Background(),
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

func (wok *Wok) ServeDir(route string, path string) error {
	ab, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	fs := http.FileServer(http.Dir(ab))
	wok.mux.Handle(route+"/", http.StripPrefix(route, fs))
	return nil
}
