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
	err := json.NewDecoder(ctx.Req.Body).Decode(&user)
	if err != nil {
		http.Error(ctx.Resp, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(user)
	if err := CreateUser(user); err != nil {
		http.Error(ctx.Resp, err.Error(), http.StatusBadRequest)
		return
	}

}

// POST: login user
func LoginHandle(ctx Context) {
	// check if username is supplied
	if ctx.Req.FormValue("username") == "" {
		ctx.Resp.Write([]byte("must include username\n"))
	}
	// check if password is supplied
	if ctx.Req.FormValue("password") == "" {
		ctx.Resp.Write([]byte("must include password\n"))
	}

	// assign username and password to vars
	username := ctx.Req.FormValue("username")
	password := ctx.Req.FormValue("password")

	// calls login function SEE: user.go for specs
	uuid, err := Login(username, password)

	if err != nil {
		ctx.Resp.Write([]byte(err.Error()))
	}

	cookie := new(http.Cookie)

	cookie.Name = "session_id"
	cookie.Value = uuid
	cookie.HttpOnly = true
	cookie.Expires = time.Now().Add(30 * time.Minute).Local()
	http.SetCookie(ctx.Resp, cookie)
	http.Redirect(ctx.Resp, ctx.Req, "/home", 200)

}

// METHOD N/A: handles the favicon, replace the favicon in the root to what you would like, default is the old gopher
func Favicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	http.ServeFile(w, r, "../favicon.ico")
}

// METHOD ANY: not found page, returns 404.html in the public dir
func NotFoundPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../public/404.html")
}

func AllUsers(ctx Context) {
	_, err := UserisAdmin(ctx)
	if err != nil {
		ctx.Resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	users, err := GetAllUsers()
	if err != nil {
		ctx.Resp.Write([]byte(err.Error()))
	}

	resp, err := json.Marshal(users)
	if err != nil {
		ctx.Resp.Write([]byte(err.Error()))
	}

	ctx.Resp.Write(resp)
}

func LogoutUser(ctx Context) {
	id, err := UserisUser(ctx)
	if err != nil {
		ctx.Resp.Write([]byte(err.Error()))
	}

	if err := Logout(id); err != nil {
		ctx.Resp.Write([]byte(http.StatusText(http.StatusBadRequest)))
	}

	cookie := new(http.Cookie)

	cookie.Name = "session_id"
	cookie.Expires = time.Now().Add(-1 * time.Second).Local()
	http.SetCookie(ctx.Resp, cookie)

}

func DeleteUserHandle(ctx Context) {
	_, err := UserisAdmin(ctx)
	if err != nil {
		ctx.Resp.Write([]byte(err.Error()))
		return
	}

	id := ctx.Req.URL.Query().Get("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		ctx.Resp.Write([]byte(http.StatusText(http.StatusInternalServerError)))
	}
	if err := DeleteUser(parsedId); err != nil {
		ctx.Resp.Write([]byte(http.StatusText(http.StatusBadRequest)))
	}

	ctx.Resp.Write([]byte(http.StatusText(http.StatusOK)))

}

func SendEmailHandle(ctx Context) {

	email := ctx.Req.URL.Query().Get("email")
	if err := SendCreateUserEmail(email); err != nil {
		ctx.Resp.Write([]byte(err.Error()))
	}

	ctx.Resp.Write([]byte("sent an email to" + email))
}

func EnqueueEmail(ctx Context) {
	email := new(Email)
	email.Address = ctx.Req.URL.Query().Get("address")
	if email.Address == "" {
		ctx.Resp.Write([]byte("must use an email address"))
	}
	if err := AddEmailtoQueue(email); err != nil {
		ctx.Resp.Write([]byte(err.Error()))
	}
	ctx.Resp.Write([]byte("added email to queue"))
}

func DequeueEmail(ctx Context) {

}

func AllEmails(ctx Context) {
	data, err := EmailsinQueue()
	if err != nil {
		ctx.Resp.Write([]byte(err.Error()))
	}

	ctx.Resp.Write(data)
}
