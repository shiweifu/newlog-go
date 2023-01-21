package main

import (
	"fmt"
	"os"
	"path"

	"github.com/urfave/cli"
	_ "github.com/urfave/cli/v2"
	_ "github.com/urfave/cli/v2/altsrc"
)

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

func NewCliApp() *cli.App {
	cliApp := cli.NewApp()
	cliApp.Name = "newlog"
	cliApp.Usage = "a simple blog system"
	cliApp.Action = func(*cli.Context) error {
		fmt.Println("this is a cli app ")
		return nil
	}

	// add action
	cliApp.Commands = []cli.Command{
		{
			Name:    "serve",
			Aliases: []string{"s"},
			Usage:   "start the server",
			Action: func(c *cli.Context) error {
				NewServer()
				return nil
			},
		},
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "init blog data folder",
			Action: func(c *cli.Context) error {
				// 读取命令行参数
				// blogPath := c.String("blog-path")
				args := c.Args()

				if len(args) == 0 {
					return fmt.Errorf("blog path is required")
				}

				blogPath := args[0]

				// 判断是否是有效的目录
				fi, err := os.Stat(blogPath)
				if err != nil || !fi.IsDir() {
					return fmt.Errorf("blog path is not a directory")
				}

				// 判断该目录是否为空
				entries, err := os.ReadDir(blogPath)
				if err != nil {
					return err
				}

				if len(entries) > 0 {
					return fmt.Errorf("blog path must not contain any files")
				}

				currDir, _ := os.Getwd()
				cfgPath := path.Join(currDir, "config.yml")
				// 判断配置是否存在
				if _, err := os.Stat(cfgPath); err != nil {
					if !os.IsNotExist(err) {
						return err
					}

					// 文件不存在，创建
					file, fileErr := os.Create(cfgPath)
					if fileErr != nil {
						return fileErr
					}
					cfgContent := fmt.Sprintf(CfgDefaultContent, blogPath)
					// 写入内容
					_, writeErr := file.WriteString(cfgContent)
					if writeErr != nil {
						return writeErr
					}
					file.Close()
				}

				return NewBlogData(blogPath)
			},
		},
	}

	return cliApp
}
