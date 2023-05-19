package wok

import "net/http"

func DefaultRouter() {
	http.HandleFunc("/", NotFoundPage)
	http.HandleFunc("/home", Homepage)
	http.HandleFunc("/user", CreatUserHandle)
	http.HandleFunc("/login", LoginHandle)
	http.HandleFunc("/favicon.ico", Favicon)
	http.HandleFunc("/all", AllUsers)
	http.HandleFunc("/logout", LogoutUser)

}
