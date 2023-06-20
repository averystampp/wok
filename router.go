package wok

import "net/http"

// includes all the default routes for user creation, login, logout, 404, favicon, and return all users
func DefaultRouter(wok *Wok) {
	wok.mux.HandleFunc("/", NotFoundPage)       // not found page, remove if you want your index to be "/"
	wok.mux.HandleFunc("/favicon.ico", Favicon) // favicon route

	wok.mux.Handle("/user", wok.Post(CreatUserHandle)) // create a user
	wok.mux.Handle("/login", wok.Post(LoginHandle))    // login to an account
	wok.mux.Handle("/all", wok.Get(AllUsers))          // show all users currently in the database
	wok.mux.Handle("/logout", wok.Get(LogoutUser))     // logout of an account
	wok.mux.Handle("/delete", wok.Delete(DeleteUserHandle))
	wok.mux.Handle("/email", wok.Get(SendEmailHandle))
	wok.mux.Handle("/addemail", wok.Post(EnqueueEmail))
	wok.mux.Handle("/getemails", wok.Get(AllEmails))
}

// Handler func is a way to declare a function that will hold a context
// if you do not want this use http.HandleFunc()
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

func NewWok(tls bool, addr, certfile, keyfile string) *Wok {
	return &Wok{
		address:  addr,
		tls:      tls,
		certFile: certfile,
		keyFile:  keyfile,
		mux:      new(http.ServeMux),
	}
}

// enforces that the client use the POST method for the passed handler
func (wok *Wok) Post(handle Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := Context{
			Resp: w,
			Req:  r,
		}

		if ctx.Req.Method != "POST" {
			ctx.Resp.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
			return
		}

		if err := handle(ctx); err != nil {
			ctx.Resp.Write([]byte(err.Error()))
			return
		}

	}
}

// enforces that the client use the GET method for the passed handler
func (wok *Wok) Get(handle Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := Context{
			Resp: w,
			Req:  r,
		}

		if ctx.Req.Method != "GET" {
			ctx.Resp.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
			return
		}

		if err := handle(ctx); err != nil {
			ctx.Resp.Write([]byte(err.Error()))
			return
		}

	}
}

// enforce the patch method for the passed handler
func (wok *Wok) Patch(handle Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := Context{
			Resp: w,
			Req:  r,
		}

		if ctx.Req.Method != "PATCH" {
			ctx.Resp.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
			return
		}

		if err := handle(ctx); err != nil {
			ctx.Resp.Write([]byte(err.Error()))
			return
		}

	}
}

// enforce put method for the passed handler
func (wok *Wok) Put(handle Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := Context{
			Resp: w,
			Req:  r,
		}

		if ctx.Req.Method != "PUT" {
			ctx.Resp.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
			return
		}

		if err := handle(ctx); err != nil {
			ctx.Resp.Write([]byte(err.Error()))
			return
		}

	}
}

func (wok *Wok) Options(handle Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := Context{
			Resp: w,
			Req:  r,
		}

		if ctx.Req.Method != "OPTIONS" {
			ctx.Resp.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
			return
		}

		if err := handle(ctx); err != nil {
			ctx.Resp.Write([]byte(err.Error()))
			return
		}

	}
}

func (wok *Wok) Delete(handle Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := Context{
			Resp: w,
			Req:  r,
		}

		if ctx.Req.Method != "DELETE" {
			ctx.Resp.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
			return
		}

		if err := handle(ctx); err != nil {
			ctx.Resp.Write([]byte(err.Error()))
			return
		}

	}
}
