package main

import (
	"html/template"
	"newlog-go/models"
	"path"

	"github.com/kataras/iris/v12"
)

var (
	posts      = make([]*models.Post, 0)
	pages      = make([]*models.Page, 0)
	categories = make([]string, 0)

	PostRoutePath  = "/post/{title}"
	PageRoutePath  = "/page/{title}"
	IndexRoutePath = "/"
	PingRoutePath  = "/ping"
)

var currentPage *models.Page
var currentPost *models.Post

var customJS template.JS
var customCSS template.CSS

func setupRoutes(app *iris.Application) {
	favPath := path.Join(config.BlogPath, "custom/favicon.ico")
	app.Favicon(favPath)
	app.HandleDir("/assets", iris.Dir("./assets"))
	app.Use(iris.Compression)
	tmpl := iris.HTML("./views", ".html").Layout("layout.html").Reload(true)
	app.Use(func(ctx iris.Context) {
		var pageTitle string
		var navLinkTitle string

		ctx.ViewData("categories", categories)
		ctx.ViewData("pages", pages)
		routerPath := ctx.GetCurrentRoute().Path()
		title := ctx.Params().Get("title")

		if routerPath == PostRoutePath {
			currentPost, _ = getPost(title)
			pageTitle = config.BlogTitle + " - " + currentPost.Title
		} else if routerPath == PageRoutePath {
			currentPage, _ = getPage(title)
			pageTitle = config.BlogTitle + " - " + currentPage.Title
			navLinkTitle = title
		} else {
			pageTitle = config.BlogTitle
		}
		ctx.ViewData("pageTitle", pageTitle)
		ctx.ViewData("blogTitle", config.BlogTitle)
		ctx.ViewData("navLinkTitle", navLinkTitle)
		ctx.Next()
	})

	tmpl.AddFunc("postsInCategory",
		func(category string) []*models.Post {
			results := make([]*models.Post, 0)
			for _, p := range posts {
				if p.Category() == category {
					results = append(results, p)
				}
			}
			return results
		},
	)

	tmpl.AddLayoutFunc("NavTitle", func(navLinkTitle string) string {
		return navLinkTitle
	})

	tmpl.AddFunc(
		"customCSS",
		func() template.CSS {
			return customCSS
		},
	)

	tmpl.AddFunc(
		"customJS",
		func() template.JS {
			return customJS
		},
	)

	app.RegisterView(tmpl)

	// 配置路由
	app.Get(PingRoutePath, pingHandler)
	app.Get(IndexRoutePath, indexHandler)
	app.Get(PostRoutePath, postHandler)
	app.Get(PageRoutePath, pageHandler)
}

func pingHandler(ctx iris.Context) {
	ctx.WriteString("let's go!")
}

func indexHandler(ctx iris.Context) {
	ctx.View("index.html")
}

func postHandler(ctx iris.Context) {
	if currentPost == nil {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	ctx.ViewData("post", currentPost)
	ctx.View("post")
}

func pageHandler(ctx iris.Context) {
	if currentPage == nil {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	ctx.ViewData("page", currentPage)
	ctx.View("page")
}
