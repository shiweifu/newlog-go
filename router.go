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

var currentPage *models.Page
var currentPost *models.Post

func setupRoutes(app *iris.Application) {
	app.HandleDir("/assets", iris.Dir("./assets"))
	app.Use(iris.Compression)
	tmpl := iris.HTML("./views", ".html").Layout("layout.html").Reload(true)
	app.Use(func(ctx iris.Context) {
		var pageTitle string

		ctx.ViewData("categories", categories)
		ctx.ViewData("pages", pages)
		routerPath := ctx.GetCurrentRoute().Path()
		title := ctx.Params().Get("title")

		if routerPath == "/post/{title}" {
			currentPost, _ = getPost(title)
			pageTitle = config.BlogTitle + " - " + currentPost.Title
		} else if routerPath == "/page/{title}" {
			currentPage, _ = getPage(title)
			pageTitle = config.BlogTitle + " - " + currentPage.Title
		} else {
			pageTitle = config.BlogTitle
		}
		ctx.ViewData("pageTitle", pageTitle)
		ctx.ViewData("blogTitle", config.BlogTitle)
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
	app.Get("/ping", ping)
	app.Get("/", index)
	app.Get("/post/{title}", post)
	app.Get("/page/{title}", page)
}

func ping(ctx iris.Context) {
	ctx.WriteString("let's go!")
}

func index(ctx iris.Context) {
	ctx.View("index.html")
}

func post(ctx iris.Context) {
	if currentPost == nil {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	ctx.ViewData("post", currentPost)
	ctx.View("post")
}

func page(ctx iris.Context) {
	if currentPage == nil {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	ctx.ViewData("page", currentPage)
	ctx.View("page")
}
