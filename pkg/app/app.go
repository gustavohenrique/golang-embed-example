package app

import (
	"log"

	"example/pkg/ui"
)

type Application struct {
	HttpServer *HttpServer
	htmlPage   *ui.HtmlPage
}

func NewApplication() *Application {
	return &Application{
		htmlPage: ui.NewHtmlPage(),
	}
}

func (a *Application) Configure() {
	httpServer := NewHttpServer(a.htmlPage)
	httpServer.Configure()
	a.HttpServer = httpServer
}

func (a *Application) Start() {
	err := a.HttpServer.Start()
	if err != nil {
		log.Fatalln("HTTP server:", err)
	}
}
