package models

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
)

type PostFrontMatter struct {
	Title        string `yaml:"title"`
	CreatedAtStr string `yaml:"created_at"`
	Private      bool   `yaml:"private"`
	Category     string `yaml:"category"`
}

func (fm *PostFrontMatter) createdAtTime() time.Time {
	// 转换时间
	var result time.Time
	if fm.CreatedAtStr != "" {
		result, _ = time.Parse("2006-01-02", fm.CreatedAtStr)
	} else {
		result = time.Now()
	}
	return result
}

func (fm *PostFrontMatter) category() string {
	if fm.Category == "" {
		return fm.createdAtTime().Format("2006")
	}
	return fm.Category
}

type Post struct {
	PostFrontMatter
	MdContent   string
	HtmlContent template.HTML
}

func NewPostFormPath(filePath, defaultCategory string) (*Post, error) {
	mdBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	matter := PostFrontMatter{}
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
		if matter.Category == "" {
			matter.Category = defaultCategory
		}
	}

	return NewPost(&matter, string(contentBytes)), nil
}

func NewPost(frontMatter *PostFrontMatter, mdContent string) *Post {
	var htmlBuff bytes.Buffer
	// 解析 markdown 文件
	if err := goldmark.Convert([]byte(mdContent), &htmlBuff); err != nil {
		panic(err)
	}
	return &Post{
		PostFrontMatter: *frontMatter,
		MdContent:       mdContent,
		HtmlContent:     template.HTML(htmlBuff.String()),
	}
}

func (p *Post) CreatedAt() time.Time {
	return p.PostFrontMatter.createdAtTime()
}

func (p *Post) CreatedAtStr() string {
	return p.CreatedAt().Format("2006-01-02")
}

func (p *Post) Path() string {
	return fmt.Sprintf("/post/%s", p.Title)
}

func (p *Post) Category() string {
	return p.PostFrontMatter.category()
}
