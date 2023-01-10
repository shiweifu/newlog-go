package main

import (
	"fmt"
	"io/ioutil"
	"newlog-go/models"
	"os"

	_ "github.com/adrg/frontmatter"
	"github.com/kataras/iris/v12"
	"gopkg.in/yaml.v2"
)

type Config struct {
	BlogPath string `yaml:"blog_path"`
}

var posts = make([]*models.Post, 0)
var categories = make([]string, 0)

func getCategories() []string {
	if len(categories) > 0 {
		return categories
	}
	results := make([]string, 0)
	for _, p := range posts {
		results = append(results, p.Category)
	}
	return results
}

func getPost(title string) (*models.Post, error) {
	for _, p := range posts {
		if p.Title == title {
			return p, nil
		}
	}
	return nil, fmt.Errorf("post not found")
}

func readConf(filename string) (*Config, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	return c, err
}

// 加载日志文件
func loadData(blogPath string) {
	// 遍历文件夹，找出所有的 .md 文件
	entries, err := os.ReadDir(blogPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	for _, item := range entries {
		if !item.IsDir() {
			fmt.Println(item.Name())
			fullPath := blogPath + "/" + item.Name()
			post, err := models.NewPostFormPath(fullPath)
			// 读取文件内容
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			posts = append(posts, post)
		}
	}
}

func main() {
	// 加载配置文件 config.yml
	config, err := readConf("config.yml")

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	loadData(config.BlogPath)

	app := iris.New()
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

	app.Listen(":8080")
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
