package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/adrg/frontmatter"
	"github.com/kataras/iris/v12"
)

var config *Config
var Env string

func NewServer() {
	// 加载配置文件 config.yml
	var err error
	config, err = readConf("./config.yml")
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	LoadData(config.BlogPath)
	go watch(config.BlogPath)
	app := iris.New()
	setupRoutes(app)
	setupLogger(app)
	hostAndPort := fmt.Sprintf(":%d", config.GetPort())
	app.Listen(hostAndPort)
}

func NewBlogData(path string) error {
	// 创建目录
	pagesPath := path + "/pages"
	err := os.Mkdir(pagesPath, os.ModePerm)
	if err != nil {
		return err
	}
	// 创建页面文件
	aboutPath := pagesPath + "/about.md"
	contactPath := pagesPath + "/contact.md"
	// 创建文件
	aboutFile, err := os.Create(aboutPath)
	if err != nil {
		return err
	}
	contactFile, err := os.Create(contactPath)
	if err != nil {
		return err
	}

	postsPath := path + "/posts"
	customPath := path + "/custom"
	staticPath := path + "/static"

	// 创建文件

	// 写入内容

	return nil

}

func main() {
	app := NewCliApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
