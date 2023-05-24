package wok

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// POST: create a user and insert them into the database
func CreatUserHandle(ctx Context) {
	user := new(User)
	err := json.NewDecoder(ctx.r.Body).Decode(&user)
	if err != nil {
		http.Error(ctx.w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(user)
	if err := CreateUser(user); err != nil {
		http.Error(ctx.w, err.Error(), http.StatusBadRequest)
		return
	}

}

// POST: login user
func LoginHandle(ctx Context) {
	// check if username is supplied
	if ctx.r.FormValue("username") == "" {
		ctx.w.Write([]byte("must include username\n"))
	}
	// check if password is supplied
	if ctx.r.FormValue("password") == "" {
		ctx.w.Write([]byte("must include password\n"))
	}

	// assign username and password to vars
	username := ctx.r.FormValue("username")
	password := ctx.r.FormValue("password")

	// calls login function SEE: user.go for specs
	uuid, err := Login(username, password)

	if err != nil {
		ctx.w.Write([]byte(err.Error()))
	}

	cookie := new(http.Cookie)

	cookie.Name = "session_id"
	cookie.Value = uuid
	cookie.HttpOnly = true
	cookie.Expires = time.Now().Add(30 * time.Minute).Local()
	http.SetCookie(ctx.w, cookie)

}

// METHOD N/A: handles the favicon, replace the favicon in the root to what you would like, default is the old gopher
func Favicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	http.ServeFile(w, r, "../favicon.ico")
}

// METHOD ANY: not found page, returns 404.html in the public dir
func NotFoundPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	http.ServeFile(w, r, "../public/404.html")
}

func AllUsers(ctx Context) {

	if err := UserisUser(ctx); err != nil {
		ctx.w.Write([]byte(err.Error()))
		return
	}

	qs := "SELECT * FROM users"
	rows, err := database.Query(qs)
	if err != nil {
		fmt.Println(err)
	}

	var users []User
	var user User
	for rows.Next() {
		rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Password, &user.Role, &user.SessionID, &user.Logged_in)
		users = append(users, user)
	}
	resp, err := json.Marshal(users)
	if err != nil {
		ctx.w.Write([]byte(err.Error()))
	}
	ctx.w.Write(resp)
}

func LogoutUser(ctx Context) {
	id := ctx.r.URL.Query().Get("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		ctx.w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
	}
	if err := Logout(parsedId); err != nil {
		ctx.w.Write([]byte(http.StatusText(http.StatusBadRequest)))
	}

	cookie := new(http.Cookie)

	cookie.Name = "session_id"
	cookie.Expires = time.Now().Add(-1 * time.Second).Local()
	http.SetCookie(ctx.w, cookie)

}
