package main

import (
	"fmt"
	"html/template"
	"log"
	"newlog-go/models"
	"os"
	"sort"

	"github.com/fsnotify/fsnotify"
)

var CfgDefaultContent string = `
blog_path: "%s"
blog_title: newlog
port: 8080
env: production
`

var AboutPageContent string = `
---
title: 关于
index: 1
private: false
---

关于本日志系统

所有在 pages 目录下的内容，将被识别为页面

pages 内容支持三个元属性：

 - title: 页面标题
 - index: 页面在导航栏中的位置
 - private: 页面是被渲染
`

var ContactPageContent string = `
---
title: 联系
index: 2
private: false
---

联系我以及其他链接
`

var PostContent string = `
---
title: 《小强升职记》读后感
private: true
created_at: "2022-11-10"
category: "2022"
---

这是一篇测试文章。

日志放在 posts 目录下，文件名为文章标题，文件内容为文章内容。

posts 目录下的文件支持四个元属性：

 - title: 文章标题。如果没有，将使用文件名作为标题
 - private: 文章是否被渲染。默认为 false
 - created_at: 文章创建时间。如果没有，将使用文件创建时间
 - category: 文章分类。如果不指定：
  - 寻找日志是否存在一级目录，如果存在，将使用一级目录作为分类
  - 如果不存在，则使用创建日期的年份作为分类
`

func contains(results []string, s string) bool {
	for _, v := range results {
		if v == s {
			return true
		}
	}
	return false
}

func getPost(title string) (*models.Post, error) {
	for _, p := range posts {
		if p.Title == title {
			return p, nil
		}
	}
	return nil, fmt.Errorf("post not found")
}

func getPage(title string) (*models.Page, error) {
	for _, p := range pages {
		if p.Title == title {
			return p, nil
		}
	}
	return nil, fmt.Errorf("page not found")
}

func getFolderFiles(path string) []string {
	entries, err := os.ReadDir(path)
	results := make([]string, 0)

	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		fileFullPath := path + "/" + entry.Name()
		if !entry.IsDir() {
			results = append(results, fileFullPath)
		}
	}
	return results
}

func appendPost(postPath, defaultCategory string) {
	post, err := models.NewPostFormPath(postPath, defaultCategory)
	if err != nil {
		log.Println(err)
		return
	}
	posts = append(posts, post)
}

func appendPage(pagePath string) {
	page, err := models.NewPageFormPath(pagePath)
	if err != nil {
		log.Println(err)
		return
	}
	pages = append(pages, page)
}

func loadPosts(blogPath string) {
	posts = make([]*models.Post, 0)
	// 遍历文件夹，找出所有的 .md 文件
	entries, err := os.ReadDir(blogPath)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	for _, entry := range entries {
		// 目录名称作为分类名
		if entry.IsDir() {
			// 读取目录下的所有文件
			fullDirPath := blogPath + "/" + entry.Name()
			mdFiles := getFolderFiles(fullDirPath)
			for _, mdFile := range mdFiles {
				appendPost(mdFile, entry.Name())
			}
		} else {
			// 文件作为文章
			appendPost(blogPath+"/"+entry.Name(), "")
		}
	}
}

func loadPages(pagesPath string) {
	// 遍历页面文件夹
	pages = make([]*models.Page, 0)
	mdFiles := getFolderFiles(pagesPath)

	for _, mdFile := range mdFiles {
		appendPage(mdFile)
	}

	// 排序
	sort.Sort(models.Pages(pages))
}

func loadCategories() {
	categories = make([]string, 0)
	if len(categories) > 0 {
		return
	}
	results := make([]string, 0)
	for _, p := range posts {
		// 如果已经包含分类，则继续
		if contains(results, p.Category()) {
			continue
		}
		results = append(results, p.Category())
	}
	categories = results
}

func loadCustom() {
	// 读取自定义的 CSS
	cssFilePath := config.BlogPath + "custom/custom.css"
	jsFilePath := config.BlogPath + "custom/custom.js"

	// 读取自定义的 JS
	cssContent, err := os.ReadFile(cssFilePath)

	if err != nil {
		log.Println(err)
	}
	jsContent, err := os.ReadFile(jsFilePath)
	if err != nil {
		log.Println(err)
	}

	customCSS = template.CSS(cssContent)
	customJS = template.JS(jsContent)
}

func watch(blogPath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// 重新加载数据
				LoadData(blogPath)
				_ = event
			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
			}
		}
	}()
	// 监听当前目录
	err = watcher.Add(blogPath)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

// 加载博客文件目录
func LoadData(blogPath string) {
	postsPath := blogPath + "/posts"
	pagesPath := blogPath + "/pages"
	loadPosts(postsPath)
	loadCategories()
	loadPages(pagesPath)
	loadCustom()
}
