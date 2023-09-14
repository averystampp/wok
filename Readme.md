### Wok.

Wok is a simple http server written in Go. Wok aims to help Gophers create simple applications without
the need overengineer their projects. Wok doesn't have any fancy features or complicated routing procedures. Its pure Go and doesn't require CGO to be enabled.

I originally set out to make a multi featured http server. Wok originally was a way to package multiple 
services for developers to configure out of the box. But, Wok feels right when its as simple as possible.

Wok is my passion project. As time goes on I will add and remove features as I need them. If you feel
like something should be added or changed in Wok feel free to make a PR. But I cannot promise that every feature
will be added.

Wok is not, nor will it probably every be as great as things like Fiber and other projects like it. But I wanted something to build and call my own when I make applications, even if that means I can't achieve the bleeding edge.

Currently Wok is under development and I will commit breaking changes without warning.

#### Simple Usage:
```
package main

import (
	"github.com/averystampp/wok"
)

func main() {
	
	app := wok.NewWok(":5000")
	
	app.Get("/", HelloFromWok)

	app.StartWok()
}

func HelloFromWok(ctx wok.Context) error {
	ctx.SendString("Hello from Wok!")
	return nil
}

```

#### Use a database (Wok only supports postgres):
```
package main

import (
	"github.com/averystampp/wok"
)

func main() {
	app := wok.NewWok(":5000")
	
	conf := &wok.Config{
		Host:            "{YOUR_POSTGRES_HOST}",
		Port:            5432,
		User:            "{YOUR_POSTGRES_USER}",
		Password:        "{YOUR_POSTGRES_PASSWORD}",
		Dbname:          "{DATABASE_NAME}",
		MigrationFolder: "./migrations",
	}
	app.WithDatabase(conf)

	app.Get("/", HelloFromWok)
	app.StartWok(dbconfig)

}

func HelloFromWok(ctx wok.Context) error {
	ctx.SendString("Hello from Wok!")
	return nil
}
```