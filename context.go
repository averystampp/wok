package wok

import (
	"database/sql"
	"fmt"
	"net/http"
)

// Context controls ResponseWriter and pointer to Request, used to extend methods
type Context struct {
	Resp http.ResponseWriter
	Req  *http.Request
	db   *sql.DB
}

func (ctx *Context) Database() (*sql.DB, error) {
	if ctx.db != nil {
		return ctx.db, nil
	}
	return nil, fmt.Errorf("no database initialized")
}
