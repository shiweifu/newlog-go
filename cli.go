package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	_ "github.com/urfave/cli/v2"
	_ "github.com/urfave/cli/v2/altsrc"
)

func NewCliApp() *cli.App {
	cliApp := cli.NewApp()
	cliApp.Name = "newlog"
	cliApp.Usage = "a simple blog system"
	cliApp.Action = func(*cli.Context) error {
		fmt.Println("boom! I say!")
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

				return NewBlogData(blogPath)
			},
		},
	}

	return cliApp
}
