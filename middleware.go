package wok

func ProtectWithToken(h Handler) Handler {
	return func(ctx Context) error {
		if err := ctx.ValidateToken(TokenSecret); err != nil {
			return err
		}
		return h(ctx)
	}
}
