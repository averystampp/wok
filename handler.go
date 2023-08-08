package wok

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// POST: create a user and insert them into the database
func CreatUserHandle(ctx Context) error {
	if err := UserisAdmin(ctx); err != nil {
		return err
	}
	user := User{}
	err := json.NewDecoder(ctx.Req.Body).Decode(&user)
	if err != nil {
		return err
	}
	defer ctx.Req.Body.Close()
	if err := CreateUser(&user); err != nil {
		return err
	}

	return nil

}

// POST: Takes in username and password from form, checks if they are blank, returns error if either are blank.
// Calls Login function which will
func LoginHandle(ctx Context) error {
	// check if username is supplied
	if ctx.Req.FormValue("username") == "" {
		ctx.Resp.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("must include username")
	}
	// check if password is supplied
	if ctx.Req.FormValue("password") == "" {
		return fmt.Errorf("must include password")
	}

	username := ctx.Req.FormValue("username")
	password := ctx.Req.FormValue("password")

	uuid, err := Login(username, password)
	if err != nil {
		return err
	}

	// session_id cookie is sent back to client to be saved and used in future requests to the server
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    uuid,
		HttpOnly: true,
		Expires:  time.Now().Add(30 * time.Minute).Local(),
	}

	http.SetCookie(ctx.Resp, cookie)

	// logged_in cookie is sent back to client to be used to handle frontend views pages for changing state
	c := &http.Cookie{
		Name:     "logged_in",
		Value:    "true",
		HttpOnly: true,
		Expires:  time.Now().Add(30 * time.Minute).Local(),
	}

	http.SetCookie(ctx.Resp, c)

	http.Redirect(ctx.Resp, ctx.Req, "/", http.StatusSeeOther)
	return nil
}

func AllUsers(ctx Context) error {
	if err := UserisAdmin(ctx); err != nil {
		return err
	}
	users, err := GetAllUsers()
	if err != nil {
		return err
	}

	if err := ctx.JSON(users); err != nil {
		return err
	}

	return nil
}

func LogoutUser(ctx Context) error {
	if err := UserisValid(ctx); err != nil {
		return err
	}

	id, err := ctx.Req.Cookie("session_id")
	if err != nil {
		return err
	}

	if err := Logout(id.Value); err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Expires:  time.Now().Add(-1 * time.Second).Local(),
		HttpOnly: true,
	}

	http.SetCookie(ctx.Resp, cookie)

	c := &http.Cookie{
		Name:     "logged_in",
		Value:    "false",
		HttpOnly: true,
		Expires:  time.Now().Add(30 * time.Minute).Local(),
	}
	http.SetCookie(ctx.Resp, c)

	http.Redirect(ctx.Resp, ctx.Req, "/", http.StatusSeeOther)
	return nil
}

func DeleteUserHandle(ctx Context) error {
	if err := UserisAdmin(ctx); err != nil {
		return err
	}
	id := ctx.Req.URL.Query().Get("id")

	if id == "" {
		return fmt.Errorf("must have an id in request params")
	}

	parsedId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	if err := DeleteUser(parsedId); err != nil {
		return err
	}

	return nil
}
