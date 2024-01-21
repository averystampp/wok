package wok

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Wok enforces its own handler to return an error, then wraps it into an
// http handler converter
type Handler func(Context) error

type Route struct {
	handlers []Handler
	method   string
}

type Router struct {
	patterns   map[string]Route
	use_static bool
	dir        string
	CertFile   string
	KeyFile    string
	db         *sql.DB
	middleware []Handler
}

func NewRouter() *Router {
	return &Router{
		patterns: make(map[string]Route),
	}
}

func (r *Router) WithDatabase(opts *DatabaseOpts) {
	Connect(opts)
}

func (r *Router) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	if r.use_static {
		if strings.HasPrefix(req.URL.Path, r.dir) {
			ctx := Context{
				Resp: wr,
				Req:  req,
			}
			entries, err := os.ReadDir("." + ctx.Req.URL.Path)
			if err != nil {
				f, err := os.Open("." + ctx.Req.URL.Path)
				if err != nil {
					ctx.Resp.Write([]byte(err.Error()))
					return
				}
				b, err := io.ReadAll(f)
				if err != nil {
					ctx.Resp.Write([]byte(err.Error()))
					return
				}
				ctx.Resp.Write(b)
				return
			}
			var bl string
			for _, entry := range entries {
				bl += fmt.Sprintf("<a href=%s>%s</a><br>", ctx.Req.URL.Path+"/"+entry.Name(), entry.Name())
			}
			ctx.Resp.Write([]byte(bl))
			return
		}
	}
	handlers, ok := r.patterns[req.URL.Path]
	if !ok {
		http.NotFound(wr, req)
		return
	}
	if req.Method != handlers.method {
		wr.WriteHeader(http.StatusMethodNotAllowed)
		wr.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	r.middleware = append(r.middleware, handlers.handlers...)
	for _, hand := range r.middleware {
		err := hand(Context{
			Resp: wr,
			Req:  req,
			db:   DB,
		})
		if err != nil {
			wr.Write([]byte(err.Error()))
			return
		}
	}
}

func (r *Router) Get(path string, h ...Handler) {
	r.patterns[path] = Route{
		method:   http.MethodGet,
		handlers: h,
	}
}

func (r *Router) Post(path string, h ...Handler) {
	r.patterns[path] = Route{
		method:   http.MethodPost,
		handlers: h,
	}
}

func (r *Router) Put(path string, h ...Handler) {
	r.patterns[path] = Route{
		method:   http.MethodPut,
		handlers: h,
	}
}

func (r *Router) Patch(path string, h ...Handler) {
	r.patterns[path] = Route{
		method:   http.MethodPatch,
		handlers: h,
	}
}

func (r *Router) Options(path string, h ...Handler) {
	r.patterns[path] = Route{
		method:   http.MethodOptions,
		handlers: h,
	}
}

func (r *Router) Head(path string, h ...Handler) {
	r.patterns[path] = Route{
		method:   http.MethodHead,
		handlers: h,
	}
}

func (r *Router) Trace(path string, h ...Handler) {
	r.patterns[path] = Route{
		method:   http.MethodTrace,
		handlers: h,
	}
}

func (r *Router) SetRootDir(root string) {
	r.dir = root
	r.use_static = true
}

var cmd = flag.String("cmd", "", "for commands to control wok")

func (r *Router) StartRouter(port string) {
	flag.Parse()

	switch *cmd {
	case "droptable":
		DropTable()
	default:
		break
	}
	fmt.Println("Server starting on port" + port)
	if r.CertFile != "" && r.KeyFile != "" {
		log.Fatal(http.ListenAndServeTLS(port, r.CertFile, r.KeyFile, r))
	} else {
		log.Fatal(http.ListenAndServe(port, r))
	}
}
