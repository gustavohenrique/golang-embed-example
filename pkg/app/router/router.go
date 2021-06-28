package router

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"

	"example/pkg/domain/entities"
	"example/pkg/ui"
)

type Router struct {
	htmlPage *ui.HtmlPage
}

func NewRouter(htmlPage *ui.HtmlPage) *Router {
	return &Router{
		htmlPage: htmlPage,
	}
}

func (r *Router) ServeStaticFiles(e *echo.Echo) {
	files := r.htmlPage.Get()
	contentHandler := echo.WrapHandler(http.FileServer(http.FS(files)))
	contentRewrite := middleware.Rewrite(map[string]string{"static/*": "/html/$1"})
	e.GET("static/*", contentHandler, contentRewrite)
}

func (r *Router) AddEndpoints(e *echo.Echo) {
	api := e.Group("api")

	// Example 1
	api.GET("/posts", func(c echo.Context) error {
		return c.Render(http.StatusOK, "posts.html", map[string]interface{}{
			"posts": getPosts(),
		})
	})

	// Example 2
	api.GET("/articles", func(c echo.Context) error {
		body, err := r.htmlPage.GetTemplate("components/posts.html")
		if err != nil {
			return c.String(500, "Body error: "+err.Error())
		}
		tmpl, err := template.New("page").Parse(body)
		if err != nil {
			return c.String(500, "Parse error: "+err.Error())
		}
		var content bytes.Buffer
		err = tmpl.Execute(&content, map[string]interface{}{"posts": getPosts()})
		if err != nil {
			return c.String(500, "Parse error: "+err.Error())
		}
		return c.HTML(200, content.String())
	})
}

func getPosts() []entities.Post {
	return []entities.Post{
		{
			Title:   "Hello World",
			Article: "It is my first post!",
		},
		{
			Title:   "Using embed",
			Article: "Embeding files in Go binary",
		},
	}
}
