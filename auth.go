package wok

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// takes in context and parses the session_id cookie, if the cookie is not present
// function returns an empty string and an error, if cookie is in the request but does
// but the user does not have the role of "user" or "admin", return an error stating they are not authorized
// if the cookie is in the request and the users role is user or admin, returns the session_id as a string and nil
// for the error
func UserisValid(ctx Context) error {
	id, err := ctx.Req.Cookie("session_id")
	if err != nil {
		return err
	}

	uuid, err := uuid.Parse(id.Value)
	if err != nil {
		return err
	}

	qs := "SELECT session_id FROM users WHERE session_id=$1"
	row := Database.QueryRow(qs, uuid.String())
	var session string
	if err := row.Scan(&session); err != nil {
		return err
	}

	return nil

}

func UserisAdmin(ctx Context) error {
	id, err := ctx.Req.Cookie("session_id")
	if err != nil {
		return err
	}

	uuid, err := uuid.Parse(id.Value)
	if err != nil {
		return err
	}

	qs := "SELECT role FROM users WHERE session_id=$1"
	row := Database.QueryRow(qs, uuid.String())
	var role string

	if err := row.Scan(&role); err != nil {
		return err
	}

	if role != "admin" {
		return fmt.Errorf(http.StatusText(http.StatusUnauthorized))
	}

	return nil

}

func csrfProtecter(ctx Context) error {
	csrfFromForm := ctx.Req.FormValue("_csrf")
	row := Database.QueryRow("SELECT * FROM csrf WHERE token=$1", csrfFromForm)
	var id int
	var token string
	var expiry string
	row.Scan(&id, &token, &expiry)

	if !hmac.Equal([]byte(token), []byte(csrfFromForm)) {
		return fmt.Errorf("crsf token is not valid")
	}

	return nil
}

// Insert this into your handlers (GET methods) to create a CSRF token for the client if they don't have one
// or if the one they have is expired, works in tandem with CsrfProtect for POST handlers

func createCsrfToken(ctx Context) (*http.Cookie, error) {
	secret := "disasecrect"
	source := rand.NewSource(time.Now().UnixNano())

	seed := strconv.Itoa(int(source.Int63()))
	salt := []byte(seed + secret)
	buf := sha256.New().Sum(salt)

	hash := fmt.Sprintf("%x", buf)
	exp := time.Now().UTC().Add(time.Minute * 30)
	_, err := Database.Exec("INSERT into csrf (token,expires) VALUES ($1,$2)", hash, exp)
	if err != nil {
		return nil, err
	}
	cookie := &http.Cookie{
		Name:     "_csrf",
		Value:    hash,
		Expires:  exp,
		HttpOnly: true,
	}
	return cookie, nil
}
