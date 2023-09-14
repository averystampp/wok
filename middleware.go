package wok

import "net/http"

// ProtectWithToken is a middleware function that will enforce the incoming request
// to a protected route has a valid api key that was set with your token secret and is not
// expired.
func ProtectWithToken(h Handler) Handler {
	return func(ctx Context) error {
		if err := ctx.ValidateToken(TokenSecret); err != nil {
			ctx.Resp.WriteHeader(http.StatusUnauthorized)
			return err
		}

		return h(ctx)
	}
}
