package wok

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

// POST: create a user and insert them into the database
func CreatUserHandle(ctx Context) error {
	user := new(User)
	err := json.NewDecoder(ctx.Req.Body).Decode(&user)
	if err != nil {
		return err
	}
	fmt.Println(user)
	if err := CreateUser(user); err != nil {
		return err
	}

	return nil

}

// POST: login user
func LoginHandle(ctx Context) error {
	// check if username is supplied
	if ctx.Req.FormValue("username") == "" {
		return fmt.Errorf("must include username")
	}
	// check if password is supplied
	if ctx.Req.FormValue("password") == "" {
		return fmt.Errorf("must include password")
	}

	// assign username and password to vars
	username := ctx.Req.FormValue("username")
	password := ctx.Req.FormValue("password")

	// calls login function SEE: user.go for specs
	uuid, err := Login(username, password)
	if err != nil {
		return err
	}

	cookie := new(http.Cookie)

	cookie.Name = "session_id"
	cookie.Value = uuid
	cookie.HttpOnly = true
	cookie.Expires = time.Now().Add(30 * time.Minute).Local()
	http.SetCookie(ctx.Resp, cookie)

	c := new(http.Cookie)
	c.Name = "logged_in"
	c.Value = "true"
	c.HttpOnly = true
	c.Expires = time.Now().Add(30 * time.Minute).Local()
	http.SetCookie(ctx.Resp, c)

	var host string
	if os.Getenv("prod") == "true" {
		host = "https://idkwtptda.com"
	} else {
		host = "http://localhost:8080"
	}
	http.Redirect(ctx.Resp, ctx.Req, host+"/home", http.StatusSeeOther)
	return nil
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

func AllUsers(ctx Context) error {
	_, err := UserisAdmin(ctx)
	if err != nil {
		return err
	}

	users, err := GetAllUsers()
	if err != nil {
		return err
	}

	resp, err := json.Marshal(users)
	if err != nil {
		return err
	}

	ctx.Resp.Write(resp)
	return nil
}

func LogoutUser(ctx Context) error {
	id, err := UserisUser(ctx)
	if err != nil {
		return err
	}

	if err := Logout(id); err != nil {
		return err
	}

	cookie := new(http.Cookie)

	cookie.Name = "session_id"
	cookie.Expires = time.Now().Add(-1 * time.Second).Local()
	http.SetCookie(ctx.Resp, cookie)

	c := new(http.Cookie)
	c.Name = "logged_in"
	c.Value = "false"
	c.HttpOnly = true
	c.Expires = time.Now().Add(30 * time.Minute).Local()
	http.SetCookie(ctx.Resp, c)
	var host string

	if os.Getenv("prod") == "true" {
		host = "https://idkwtptda.com"
	} else {
		host = "http://localhost:8080"
	}
	http.Redirect(ctx.Resp, ctx.Req, host+"/home", http.StatusSeeOther)
	return nil
}

func DeleteUserHandle(ctx Context) error {
	_, err := UserisAdmin(ctx)
	if err != nil {
		return err
	}

	id := ctx.Req.URL.Query().Get("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	if err := DeleteUser(parsedId); err != nil {
		return err
	}

	return nil
}

func SendEmailHandle(ctx Context) error {

	email := ctx.Req.URL.Query().Get("email")
	if err := SendCreateUserEmail(email); err != nil {
		return err
	}
	return nil
}

func EnqueueEmail(ctx Context) error {
	email := new(Email)
	email.Address = ctx.Req.URL.Query().Get("address")
	if email.Address == "" {
		return fmt.Errorf("must use an email address")
	}
	if err := AddEmailtoQueue(email); err != nil {
		return err
	}
	ctx.Resp.Write([]byte("added email to queue"))
	return nil
}

func DequeueEmail(ctx Context) error {
	return nil
}

func AllEmails(ctx Context) error {
	data, err := EmailsinQueue()
	if err != nil {
		return err
	}

	ctx.Resp.Write(data)
	return nil
}
