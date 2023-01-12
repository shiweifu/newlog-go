package main

import (
	"newlog-go/models"

	"github.com/kataras/iris/v12"
)

var (
	posts      = make([]*models.Post, 0)
	pages      = make([]*models.Page, 0)
	categories = make([]string, 0)
)

func setupRoutes(app *iris.Application) {
	app.HandleDir("/assets", iris.Dir("./assets"))
	app.Use(iris.Compression)
	tmpl := iris.HTML("./views", ".html").Layout("layout.html").Reload(true)
	app.Use(func(ctx iris.Context) {
		ctx.ViewData("pageTitle", "hello")
		ctx.ViewData("categories", categories)
		ctx.ViewData("pages", pages)
		ctx.Next()
	})

	tmpl.AddFunc("postsInCategory", func(category string) []*models.Post {
		results := make([]*models.Post, 0)
		for _, p := range posts {
			if p.Category() == category {
				results = append(results, p)
			}
		}
		return results
	})

	app.RegisterView(tmpl)

	// 配置路由
	// GET: http://localhost:8080/hello
	app.Get("/ping", ping)
	app.Get("/html", index)
	app.Get("/post/{title}", post)
	app.Get("/page/{title}", page)
}

func ping(ctx iris.Context) {
	ctx.WriteString("Hello from the server!")
}

func index(ctx iris.Context) {
	ctx.View("index.html")
}

func post(ctx iris.Context) {
	title := ctx.Params().Get("title")
	var post *models.Post
	post, err := getPost(title)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	ctx.ViewData("post", post)
	ctx.View("post")
}

func page(ctx iris.Context) {
	title := ctx.Params().Get("title")
	var page *models.Page
	page, err := getPage(title)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}

	ctx.ViewData("page", page)
	ctx.View("page")
}
