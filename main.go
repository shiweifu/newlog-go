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

	app.RegisterView(tmpl)

	// 配置路由
	// GET: http://localhost:8080/hello
	app.Get("/ping", ping)
	app.Get("/html", getHtmlTmp)

	app.Listen(":8080")
}

func ping(ctx iris.Context) {
	ctx.WriteString("Hello from the server!")
}

func getHtmlTmp(ctx iris.Context) {
	ctx.ViewData("title", "My Page Title")
	ctx.View("index.html")
}
