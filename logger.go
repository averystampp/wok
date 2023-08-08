package wok

import (
	"fmt"
	"log"
	"os"
)

const (
	WokErr  string = "[ERROR]"
	WokWarn string = "[WARNING]"
	WokInfo string = "[INFO]"
)

type Log struct {
	WokLogger
	file *os.File
}

type WokLogger interface {
	Info(Context, string)
	Warn(Context, string)
	Error(Context, string)
}

func (l *Log) NewLogFile(filename string) error {
	f, err := os.OpenFile(filename+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	l.file = f
	return nil
}

func (l *Log) Info(ctx Context, d string) {
	e := fmt.Sprintf(WokInfo+" MSG: %s %s", ctx.Req.URL, d)
	log.SetOutput(l.file)
	log.Println(e)
	defer l.file.Close()
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
