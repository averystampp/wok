package wok

import "net/http"

// includes all the default routes for user creation, login, logout, 404, favicon, and return all users
func DefaultRouter() {
	http.HandleFunc("/", NotFoundPage)
	http.HandleFunc("/favicon.ico", Favicon)

	http.Handle("/user", Post(CreatUserHandle))
	http.Handle("/login", Post(LoginHandle))
	http.Handle("/all", Get(AllUsers))
	http.Handle("/logout", Get(LogoutUser))
}

// Handler func is a way to declare a function that will hold a context
// if you do not want this use http.HandleFunc()
type Handler func(Context)

type Context struct {
	w http.ResponseWriter
	r *http.Request
}

func Post(handle Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
			return
		}

		ctx := Context{
			w: w,
			r: r,
		}

		handle(ctx)
	}
}

func Get(handle Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
			return
		}

		ctx := Context{
			w: w,
			r: r,
		}

		handle(ctx)
	}
}
