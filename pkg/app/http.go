package app

import (
	"context"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	r "example/pkg/app/router"
	"example/pkg/ui"
)

type HttpServer struct {
	rawServer *echo.Echo
	router    *r.Router
	htmlPage  *ui.HtmlPage
}

type TemplateRenderer struct {
	templates *template.Template
}

func NewHttpServer(htmlPage *ui.HtmlPage) *HttpServer {
	rawServer := echo.New()
	rawServer.HideBanner = true
	rawServer.Debug = true

	return &HttpServer{
		rawServer: rawServer,
		router:    r.NewRouter(htmlPage),
		htmlPage:  htmlPage,
	}
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func (s *HttpServer) GetRawServer() *echo.Echo {
	return s.rawServer
}

func (s *HttpServer) Configure() {
	s.addMiddlewares()

	e := s.rawServer
	t, err := template.ParseFS(s.htmlPage.Get(), "*/*.html", "*/*/*.html")
	if err != nil {
		log.Fatalln("Template parse error:", err, "t=", t)
	}
	e.Renderer = &TemplateRenderer{
		templates: t,
	}

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok!")
	})

	s.router.ServeStaticFiles(e)
	s.router.AddEndpoints(e)
}

func (s *HttpServer) addMiddlewares() {
	e := s.rawServer
	// e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowMethods: []string{
			http.MethodOptions,
			http.MethodGet,
			http.MethodPut,
			http.MethodPost,
			http.MethodDelete,
		},
		ExposeHeaders: []string{
			"grpc-status",
			"grpc-message",
		},
		AllowHeaders: []string{
			"Accept",
			"Accept-Encoding",
			"Authorization",
			"Content-Type",
			"Content-Length",
			"grpc-status",
			"grpc-message",
			"Host",
			"User-Agent",
			"XMLHttpRequest",
			"X-Requested-With",
			"X-Request-ID",
			"X-CSRF-Token",
			"x-user-agent",
			"x-grpc-web",
			"X-Amzn-Trace-Id",
			"X-Forwarded-For",
			"X-Forwarded-Port",
			"X-Real-Ip",
		},
	}))
	e.Use(middleware.BodyLimit("5M"))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
}

func (s *HttpServer) Start() error {
	port := ":8181"
	e := s.rawServer
	go func() {
		log.Fatal(e.Start(port))
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGQUIT)
	<-quit
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
