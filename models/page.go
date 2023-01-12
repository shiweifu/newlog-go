package models

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
)

type PageFrontMatter struct {
	Title   string `yaml:"title"`
	Private bool   `yaml:"private"`
	Index   int    `yaml:"index"`
}

type Page struct {
	PageFrontMatter
	MdContent   string
	HtmlContent template.HTML
}

func NewPageFormPath(filePath string) (*Page, error) {
	mdBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	matter := PageFrontMatter{}
	contentBytes, err := frontmatter.Parse(bytes.NewReader(mdBytes), &matter)
	if err != nil {
		return nil, err
	}

	if matter.Title == "" {
		// last path of file path
		// 得到文件名称没有扩展名
		fullName := strings.Split(filePath, "/")[len(strings.Split(filePath, "/"))-1]
		fileName := strings.Split(fullName, ".")[0]
		matter.Title = fileName
	}

	return NewPage(&matter, string(contentBytes)), nil
}

func NewPage(frontMatter *PageFrontMatter, mdContent string) *Page {
	var htmlBuff bytes.Buffer
	// 解析 markdown 文件
	if err := goldmark.Convert([]byte(mdContent), &htmlBuff); err != nil {
		panic(err)
	}
	return &Page{
		PageFrontMatter: *frontMatter,
		MdContent:       mdContent,
		HtmlContent:     template.HTML(htmlBuff.String()),
	}
}

func (p *Page) Path() string {
	return fmt.Sprintf("/page/%s", p.Title)
}
