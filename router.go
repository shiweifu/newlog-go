package main

import (
	"fmt"
	"newlog-go/models"

	"github.com/kataras/iris/v12"
)

var (
	posts      = make([]*models.Post, 0)
	categories = make([]string, 0)
)

func setupRoutes(app *iris.Application) {
	app.HandleDir("/assets", iris.Dir("./assets"))
	app.Use(iris.Compression)
	tmpl := iris.HTML("./views", ".html").Layout("layout.html").Reload(true)
	tmpl.AddLayoutFunc("pageTitle", func() string {
		return "My Page Title"
	})

	tmpl.AddFunc("postsInCategory", func(category string) []*models.Post {
		results := make([]*models.Post, 0)
		for _, p := range posts {
			fmt.Println(p.Category, category, p)
			if p.Category == category {
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

}

func ping(ctx iris.Context) {
	ctx.WriteString("Hello from the server!")
}

func index(ctx iris.Context) {
	ctx.ViewData("title", "My Page Title")
	ctx.ViewData("categories", getCategories())
	ctx.View("index.html")
}

func post(ctx iris.Context) {
	title := ctx.Params().Get("title")
	var post *models.Post
	post, err := getPost(title)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
	}
	// 取消转义HTML
	ctx.ViewData("post", post)
	ctx.View("post")
}
