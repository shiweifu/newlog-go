package main

import (
	"fmt"
	"log"
	"newlog-go/models"
	"os"

	"github.com/fsnotify/fsnotify"
)

func getCategories() []string {
	if len(categories) > 0 {
		return categories
	}
	results := make([]string, 0)
	for _, p := range posts {
		// 如果已经包含分类，则继续
		if contains(results, p.Category) {
			continue
		}
		results = append(results, p.Category)
	}
	categories = results
	return categories
}

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

// 加载日志文件
func loadData(blogPath string) {
	posts = make([]*models.Post, 0)
	categories = make([]string, 0)
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
	_ = getCategories()
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
				loadData(blogPath)
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
