package wok

import "net/http"

// includes all the default routes for user creation, login, logout, 404, favicon, and return all users
func DefaultRouter() {
	http.HandleFunc("/", NotFoundPage)
	http.HandleFunc("/user", CreatUserHandle)
	http.HandleFunc("/login", LoginHandle)
	http.HandleFunc("/favicon.ico", Favicon)
	http.HandleFunc("/all", AllUsers)
	http.HandleFunc("/logout", LogoutUser)
}
