package wok

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Token struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Expires int64  `json:"expires"`
}

var TokenSecret string

// Context controls ResponseWriter and pointer to Request, used to extend methods
type Context struct {
	Resp http.ResponseWriter
	Req  *http.Request
	Ctx  context.Context
}

// Write data to JSON
func (ctx *Context) JSON(data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	ctx.Resp.Header().Set("Content-Type", "application/json")
	ctx.Resp.Write(body)
	return nil
}

// Set key and value pairs for Ctx
func (ctx *Context) CtxSetKey(key any, val any) {
	ctx.Ctx = context.WithValue(ctx.Ctx, key, val)
}

// Returns value of key as a string
func (ctx *Context) GetCtxKey(key any) string {
	valuefromCtx := ctx.Ctx.Value(key)
	value := fmt.Sprintf("%v", valuefromCtx)
	return value
}

// Set any amount of cookies (does not make any CORS assumptions)
func (ctx *Context) SetCookie(c ...*http.Cookie) {
	for _, cookie := range c {
		http.SetCookie(ctx.Resp, cookie)
	}
}

// Get a cookie from the request by its name
func (ctx *Context) GetCookie(name string) (*http.Cookie, error) {
	return ctx.Req.Cookie(name)
}

// Sends the data as a string in the resposne
func (ctx *Context) SendString(data string) error {
	ctx.Resp.Header().Set("Content-Type", "text/plain")
	return fmt.Errorf("%s", data)
}

// Makes a request to a URL. Returns a response or an error
func (ctx *Context) MakeRequest(method, url string, data io.Reader) (*http.Response, error) {
	var client http.Client

	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Reads the request body and returns the data as an array of bytes or an error
func (ctx *Context) ReadBody() ([]byte, error) {
	return io.ReadAll(ctx.Req.Body)
}

func (ctx *Context) Redirect(url string) {
	http.Redirect(ctx.Resp, ctx.Req, url, http.StatusPermanentRedirect)
}

func (ctx *Context) CreateToken(secret string) (*http.Cookie, error) {
	t1 := Token{
		Name:    "APIKEY",
		Value:   "thisisyourkey",
		Expires: time.Now().Add(time.Minute * 30).Unix(),
	}

	key, err := json.Marshal(&t1)

	if err != nil {
		return nil, err
	}
	encoded := base64.RawURLEncoding.EncodeToString(key)

	hash := hmac.New(sha256.New, []byte(secret))

	_, err = hash.Write([]byte(encoded))
	if err != nil {
		return nil, err
	}

	cval := fmt.Sprintf("%s.%x", encoded, hash.Sum(nil))

	return &http.Cookie{
		Name:    "apikey",
		Value:   cval,
		Expires: time.Now().Add(time.Minute * 30),
	}, nil

}

func (ctx *Context) ValidateToken(secret string) error {
	hash := hmac.New(sha256.New, []byte(secret))
	cook, err := ctx.GetCookie("apikey")
	if err != nil {
		return err
	}
	vals := strings.Split(cook.Value, ".")
	b, err := base64.RawURLEncoding.DecodeString(vals[0])
	if err != nil {
		return err
	}
	t := &Token{}

	if err := json.Unmarshal(b, t); err != nil {
		return err
	}

	if t.Expires < time.Now().Unix() {
		data := map[string]string{"error": "token is expired"}
		return ctx.JSON(data)
	}

	hash.Write([]byte(vals[0]))
	sum := hash.Sum(nil)
	asString := fmt.Sprintf("%x", sum)
	if !hmac.Equal([]byte(asString), []byte(vals[1])) {
		fmt.Println(asString)
		fmt.Println(vals[1])
		dataerror := map[string]string{"error": "token does not match"}
		return ctx.JSON(dataerror)
	}
	return nil
}
