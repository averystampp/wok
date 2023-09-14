package wok

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
)

// Token is used to protect routes with the ProtectWithToken middleware.
// See middleware.go for more information
type Token struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Expires int64  `json:"expires"`
}

var TokenSecret string

// Context controls ResponseWriter and pointer to Request, used to extend methods
type Context struct {
	Ctx  context.Context
	Resp http.ResponseWriter
	Req  *http.Request
}

// JSON will take in an interface and marshal it to []byte. JSON does not
// enforce a content length and will write any amount of data to an array.
func (ctx *Context) JSON(data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	ctx.Resp.Header().Set("Content-Type", "application/json")
	ctx.Resp.Write(body)
	return nil
}

// SetCtxKey will take in any type and set it the the wok.Context.Ctx.
// Note that wok.Context(s) are pooled and short lived. Use the session to hold
// ephemeral data for a longer lived time.
func (ctx *Context) SetCtxKey(key any, val any) {
	ctx.Ctx = context.WithValue(ctx.Ctx, key, val)
}

// GetCtxKey will return the value stored at the given key (if any).
// GetCtxKey will return nil if the key is not in the given context.
func (ctx *Context) GetCtxKey(key any) any {
	return ctx.Ctx.Value(key)
}

// SetCookie will take any amount of cookies and set them to the wok.Context.ResponseWriter.
// Note that SetCookies does not make any CORS assumptions or enforce and cookies rules
// it is up to the developer to set the cookie paramater
// See https://pkg.go.dev/net/http#Cookie for documentation on cookies
func (ctx *Context) SetCookies(c ...*http.Cookie) {
	for _, cookie := range c {
		http.SetCookie(ctx.Resp, cookie)
	}
}

// GetCookie queries the request by cookie name and returns a *http.Cookie if present.
func (ctx *Context) GetCookie(name string) (*http.Cookie, error) {
	return ctx.Req.Cookie(name)
}

// SendString will return the string arguement as a string in the resposne
// body.
func (ctx *Context) SendString(data string) {
	ctx.Resp.Header().Set("Content-Type", "text/plain")
	ctx.Resp.Write([]byte(data))
}

// MakeRequest will send a request to the given URL with the given data.
// Returns a *http.Response and an error.
//
// Note that MakeRequest does not have any cookies or other data.
// It only forwards the data in the request body.
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

// Redirect will make a http.Redirect to the given url with the parent
// Context as the Response and Request.
func (ctx *Context) Redirect(url string) {
	http.Redirect(ctx.Resp, ctx.Req, url, http.StatusPermanentRedirect)
}

// CreateToken creates a apikey cookie that can then be sent back to the client
// in a response to validate future requests to routes that use the ProtectWithToken
// middleware. CreateToken uses a JWT like way of encrypting a token.
//
// CreateToken uses TokenSecret as its signing key.
func (ctx *Context) CreateToken() (*http.Cookie, error) {
	t1 := Token{
		Name:  "APIKEY",
		Value: "thisisyourkey",
	}

	key, err := json.Marshal(&t1)

	if err != nil {
		return nil, err
	}
	encoded := base64.RawURLEncoding.EncodeToString(key)

	hash := hmac.New(sha256.New, []byte(TokenSecret))

	_, err = hash.Write([]byte(encoded))
	if err != nil {
		return nil, err
	}

	cval := fmt.Sprintf("%s.%x", encoded, hash.Sum(nil))

	return &http.Cookie{
		Name:  "apikey",
		Value: cval,
	}, nil
}

// ValidateToken takes in the apikey cookie from the request and validates it based off
// your secret that you set during init. It also ensures that the token taken from the server session
// is not expired
//
// Note: If you do not set a secret, then you will not be able to use this function safely.
func (ctx *Context) ValidateToken(secret string) error {
	hash := hmac.New(sha256.New, []byte(secret))
	cook, err := ctx.GetCookie("apikey")
	if err == http.ErrNoCookie {
		return ctx.JSON(map[string]string{"error": "request is not authorized"})
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

	hash.Write([]byte(vals[0]))
	sum := hash.Sum(nil)
	asString := fmt.Sprintf("%x", sum)

	if !hmac.Equal([]byte(asString), []byte(vals[1])) {
		return fmt.Errorf("tokens do not match")
	}

	return nil
}

// MakeAuthAPICall is for making an HTTP request to a Wok.Handler that leverages the
// ProtectWithToken middleware. This adds the api key to the request and returns a
// *http.Response or an error.
// Deprecated use Forward() to push requests with authentication
func (ctx *Context) MakeAuthAPICall(method string, url string, body io.ReadCloser) (*http.Response, error) {
	client := http.Client{}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	apikey, err := ctx.Req.Cookie("apikey")
	if err != nil {
		return nil, err
	}
	req.AddCookie(apikey)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ReadData takes an io.ReadCloser and reads it into a *bytes.Buffer
// writing 256 bytes into the buffer at once.
func (ctx *Context) ReadData(data io.ReadCloser) ([]byte, error) {
	b := make([]byte, 256)
	buf := bytes.NewBuffer(nil)

	defer ctx.Req.Body.Close()

	for {
		n, err := data.Read(b)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}

		if n == 0 {
			break
		}

		_, err = buf.Write(b[:n])
		if err != nil {
			return nil, err
		}

	}
	return buf.Bytes(), nil
}

// SendHTML takes in your HTML as a string and any data you want to pass to the template
// and parses it into a *html.Template, it also sets the content type to text/html, and
// returns template.Execute with the ResponseWriter and passed data as arguments
func (ctx *Context) SendHTML(html string, data any) error {
	tmpl, err := template.New("tmp").Parse(html)
	if err != nil {
		return err
	}

	ctx.Resp.Header().Set("Content-Type", "text/html")

	return tmpl.Execute(ctx.Resp, data)
}

func (ctx *Context) SetValue(key string, val interface{}) error {
	return wokSession.AddItem(key, val)

}

func (ctx *Context) DeleteValue(key string) error {
	return wokSession.DeleteItem(key)

}

func (ctx *Context) GetValue(key string) (interface{}, error) {
	return wokSession.RetrieveItem(key)
}
