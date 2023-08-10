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
	file *os.File
}

func (l *Log) NewLogFile() error {
	f, err := os.OpenFile("wok.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	l.file = f
	return nil
}

func (l *Log) Info(ctx Context) {
	e := fmt.Sprintf(WokInfo+" %s", ctx.Req.URL)
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
