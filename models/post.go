package models

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
)

type FrontMatter struct {
	Title        string `yaml:"title"`
	CreatedAtStr string `yaml:"created_at"`
	Private      bool   `yaml:"private"`
	Category     string `yaml:"category"`
}

func (fm *FrontMatter) createdAtTime() time.Time {
	// 转换时间
	var result time.Time
	if fm.CreatedAtStr != "" {
		result, _ = time.Parse("2006-01-02", fm.CreatedAtStr)
	} else {
		result = time.Now()
	}
	return result
}

func (fm *FrontMatter) category() string {
	if fm.Category == "" {
		return fm.createdAtTime().Format("2006")
	}
	return fm.Category
}

type Post struct {
	Title       string
	CreatedAt   time.Time
	Category    string
	Private     bool
	MdContent   string
	HtmlContent string
}

func NewPostFormPath(filePath string) (*Post, error) {
	mdBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	matter := FrontMatter{}
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

	return NewPost(&matter, string(contentBytes)), nil
}

func NewPost(frontMatter *FrontMatter, mdContent string) *Post {
	var htmlBuff bytes.Buffer
	// 解析 markdown 文件
	if err := goldmark.Convert([]byte(mdContent), &htmlBuff); err != nil {
		panic(err)
	}
	fmt.Println(htmlBuff.String())
	return &Post{
		Title:       frontMatter.Title,
		CreatedAt:   frontMatter.createdAtTime(),
		Private:     frontMatter.Private,
		MdContent:   mdContent,
		HtmlContent: htmlBuff.String(),
		Category:    frontMatter.category(),
	}
}

func (p *Post) CreatedAtStr() string {
	return p.CreatedAt.Format("2006-01-02")
}
