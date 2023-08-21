package wok

import (
	"os"
	"sync"

	"golang.org/x/exp/slog"
)

type woklog struct {
	sl *slog.Logger
	mu sync.Mutex
}

func newLogger() (*woklog, error) {
	file, err := os.OpenFile("log.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	return &woklog{
		sl: slog.New(slog.NewTextHandler(file, nil)),
	}, nil
}

func (l *woklog) Info(ctx *Context, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.sl.InfoContext(ctx.Ctx, msg, "route", ctx.Req.URL.Path, "method", ctx.Req.Method)
}

func (l *woklog) Warn(ctx *Context, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.sl.WarnContext(ctx.Ctx, msg, "route", ctx.Req.URL.Path, "method", ctx.Req.Method)
}
