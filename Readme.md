### Wok.

Wok is a simple http server written in Go. Wok aims to help Gophers create simple applications without
the need overengineer their projects. Wok doesn't have any fancy features or complicated routing proceedures. Wok is under 1000 lines and works with only a few dependacies. Its pure Go and doesn't require CGO to be enabled.

I originally set out to make a multi featured http server. Wok originally was a way to package multiple 
services for developers to configure out of the box. But, Wok feels right when its as simple as possible.

Wok is my passion project. As time goes on I will add and remove features as I need them. If you feel
like something should be added or changed in Wok feel free to make a PR. But I cannot promise that every feature
will be added.

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
	

	dbconfig := wok.DbConfig{
		Host:            "localhost",
		Port:            5432,
		User:            "postgres",
		Password:        "docker",
		Dbname:          "postgres",
		MigrationFolder: "./migrations",
	}

	app.Get("/", HelloFromWok)

	app.StartWok(dbconfig)

}

func HelloFromWok(ctx wok.Context) error {
	ctx.Resp.Write([]byte("Hello from Wok!"))
	return nil
}

```