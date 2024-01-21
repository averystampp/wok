package wok

func (r *Router) With(h ...Handler) {
	r.middleware = append(r.middleware, h...)
}
