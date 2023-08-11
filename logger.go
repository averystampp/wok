package wok

import (
	"fmt"
	"log"
	"os"
	"sync"

	"golang.org/x/exp/slog"
)

const (
	WokErr  string = "[ERROR]"
	WokWarn string = "[WARNING]"
	WokInfo string = "[INFO]"
)

type Log struct {
	file *os.File
	sl   *slog.Logger
	mu   sync.Mutex
}

func NewLogger() (*Log, error) {
	file, err := os.Create("./log/log.txt")
	if err != nil {
		return nil, err
	}

	return &Log{
		sl: slog.New(slog.NewJSONHandler(file, nil)),
	}, nil
}

func (l *Log) Info(ctx Context) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.sl.Info("RESOLVE", "route", ctx.Req.URL.Path, "method", ctx.Req.Method)
}

func (l *Log) Warn(ctx Context, d string) {
	e := fmt.Sprintf(WokWarn+" %s MSG: %s", ctx.Req.URL, d)
	log.SetOutput(l.file)
	log.Println(e)
	defer l.file.Close()
}

func (l *Log) Error(ctx Context, d string) {
	e := fmt.Sprintf(WokErr+" %s MSG: %s", ctx.Req.URL, d)
	log.SetOutput(l.file)
	log.Println(e)
	defer l.file.Close()
}

func (l *Log) ReadLogFile() {
	b, err := os.ReadFile("./log/log.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(b)
}
