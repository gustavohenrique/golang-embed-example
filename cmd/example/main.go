package main

import (
	"example/pkg/app"
)

func main() {
	application := app.NewApplication()
	application.Configure()
	application.Start()
}
