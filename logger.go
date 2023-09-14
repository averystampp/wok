package wok

import (
	"log/slog"
	"os"
	"sync"
)

type WokLog struct {
	Logger *slog.Logger
	mu     sync.Mutex
}

func newLogger() (*WokLog, error) {
	file, err := os.OpenFile("log.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	return &WokLog{
		Logger: slog.New(slog.NewTextHandler(file, nil)),
	}, nil
}

func (l *WokLog) Info(ctx *Context, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.Logger.InfoContext(ctx.Ctx, msg, "route", ctx.Req.URL.Path, "method", ctx.Req.Method)
}

func (l *WokLog) General(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}
