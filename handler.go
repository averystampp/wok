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
	email.Name = ctx.Req.URL.Query().Get("name")
	if email.Name == "" {
		return fmt.Errorf("must use an email address")
	}

	if err := email.AddEmailtoQueue(); err != nil {
		return err
	}
	return nil
}

func DequeueEmail(ctx Context) error {
	_, err := UserisAdmin(ctx)
	if err != nil {
		return err
	}

	id := ctx.Req.URL.Query().Get("id")
	if id == "" {
		return fmt.Errorf("must include email id")
	}
	dbid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	email := new(Email)
	email.Id = dbid
	if err := email.RemoveEmailFromQueue(); err != nil {
		return err
	}
	return nil
}

func AllEmails(ctx Context) error {
	_, err := UserisAdmin(ctx)
	if err != nil {
		return err
	}
	data, err := EmailsinQueue()
	if err != nil {
		return err
	}

	var emails []Email

	if err := json.Unmarshal(data, &emails); err != nil {
		return err
	}

	ctx.Resp.Write([]byte(string(data)))
	return nil
}
