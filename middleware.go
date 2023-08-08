package wok

import (
	"net/http"
)

func CSRFProtect(h Handler) Handler {
	return func(ctx Context) error {
		if err := csrfProtecter(ctx); err != nil {
			return err
		}

		if err := h(ctx); err != nil {
			return err
		}
		return nil
	}
}

func CSRFCreate(h Handler) Handler {
	return func(ctx Context) error {
		cookie, err := createCsrfToken(ctx)
		if err != nil {
			return err
		}
		http.SetCookie(ctx.Resp, cookie)

		if err := h(ctx); err != nil {
			return err
		}

		return nil
	}
}
