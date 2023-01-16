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
