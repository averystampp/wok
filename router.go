package wok

import "net/http"

// includes all the default routes for user creation, login, logout, 404, favicon, and return all users
func DefaultRouter() {
	http.HandleFunc("/", NotFoundPage)       // not found page, remove if you want your index to be "/"
	http.HandleFunc("/favicon.ico", Favicon) // favicon route

	http.Handle("/user", Post(CreatUserHandle)) // create a user
	http.Handle("/login", Post(LoginHandle))    // login to an account
	http.Handle("/all", Get(AllUsers))          // show all users currently in the database
	http.Handle("/logout", Get(LogoutUser))     // logout of an account
}

// Handler func is a way to declare a function that will hold a context
// if you do not want this use http.HandleFunc()
type Handler func(Context)

type Context struct {
	w http.ResponseWriter
	r *http.Request
}

// enforces that the client use the POST method for the passed handler
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

// enforces that the client use the GET method for the passed handler
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
