package main

import (
	"fmt"
	"log"
	"os"
	"path"

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

func newBlogDataPages(basePath string) error {
	// 创建目录
	pagesPath := path.Join(basePath, "pages")
	err := os.Mkdir(pagesPath, os.ModePerm)
	if err != nil {
		return err
	}
	// 创建页面文件
	aboutPath := path.Join(pagesPath, "about.md")
	contactPath := path.Join(pagesPath, "contact.md")
	// 创建文件
	aboutFile, err := os.Create(aboutPath)
	if err != nil {
		return err
	}
	// 写入内容
	_, err = aboutFile.WriteString(AboutPageContent)
	if err != nil {
		return err
	}
	aboutFile.Close()

	contactFile, err := os.Create(contactPath)
	if err != nil {
		return err
	}
	_, err = contactFile.WriteString(ContactPageContent)
	if err != nil {
		return err
	}
	contactFile.Close()
	return nil
}

func newBlogDataPosts(basePath string) error {
	// 创建目录
	postsPath := path.Join(basePath, "posts")
	err := os.Mkdir(postsPath, os.ModePerm)
	if err != nil {
		return err
	}
	// 创建页面文件
	testPath := path.Join(postsPath, "test.md")
	// 创建文件
	testFile, err := os.Create(testPath)
	if err != nil {
		return err
	}
	// 写入内容
	_, err = testFile.WriteString(PostContent)
	if err != nil {
		return err
	}
	testFile.Close()
	return nil
}

func newBlogDataCustom(basePath string) error {
	customPath := path.Join(basePath, "custom")
	err := os.Mkdir(customPath, os.ModePerm)
	if err != nil {
		return err
	}
	cssPath := path.Join(customPath, "custom.css")
	// 创建文件
	cssFile, err := os.Create(cssPath)
	if err != nil {
		return err
	}
	cssFile.Close()

	jsPath := path.Join(customPath, "custom.js")
	// 创建文件
	jsFile, err := os.Create(jsPath)
	if err != nil {
		return err
	}
	jsFile.Close()
	return nil
}

func NewBlogData(path string) error {
	fmt.Println("Creating new blog data at " + path + "...")
	fmt.Println("Creating blog pages data...")
	err := newBlogDataPages(path)
	if err != nil {
		return err
	}

	fmt.Println("Creating blog posts data...")
	err = newBlogDataPosts(path)
	if err != nil {
		return err
	}

	fmt.Println("Creating blog custom data...")
	err = newBlogDataCustom(path)
	if err != nil {
		return err
	}
	fmt.Println("Create blog data success!")
	return nil
}

func main() {
	app := NewCliApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
