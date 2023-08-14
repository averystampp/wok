### Wok.

Wok is a simple http server written in Go. Wok aims to help Gophers create simple applications without
the need overengineer their projects. Wok doesn't have any fancy features or complicated routing proceedures. Wok is under 1000 lines and works with only a few dependacies. Its pure Go and doesn't require CGO to be enabled.

I originally set out to make a multi featured http server. Wok originally was a way to package multiple 
services for developers to configure out of the box. But, Wok feels right when its as simple as possible.

Wok is my passion project. As time goes on I will add and remove features as I need them. If you feel
like something should be added or changed in Wok feel free to make a PR. But I cannot promise that every feature
will be added.

Wok is not, nor will it probably every be as great as things like Fiber and other routers. But I wanted something
to build and call my own when I make applications, even if that means I can't achieve the bleeding edge that large OS projects produce.

Currently Wok is under development and I will commit breaking changes without warning.

#### Simple Usage:
```
package main

import (
	"github.com/averystampp/wok"
)

func main() {
	config := wok.Wok{
		Address:  ":8080",
		Database: false,
	}

	app := wok.NewWok(config)
	
	app.Get("/", HelloFromWok)

	app.StartWok()
}

func HelloFromWok(ctx wok.Context) error {
	return ctx.SendString("Hello from Wok!")
}

```

#### Use a database (Wok only supports postgres):
```
package main

import (
	"github.com/averystampp/wok"
)

func main() {
	config := wok.Wok{
		Address:  ":8080",
		Database: true,
	}

	app := wok.NewWok(config)
	
	dbconfig := wok.DbConfig{
		Host:            "{YOUR_POSTGRES_HOST}",
		Port:            5432,
		User:            "{YOUR_POSTGRES_USER}",
		Password:        "{YOUR_POSTGRES_PASSWORD}",
		Dbname:          "{DATABASE_NAME}",
		MigrationFolder: "./migrations",
	}

	app.Get("/", HelloFromWok)

	app.StartWok(dbconfig)

}

func HelloFromWok(ctx wok.Context) error {
	return ctx.SendString("Hello from Wok!")
}
```