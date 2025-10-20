package main

import (
	"os"
	"fmt"
	"log"
	"path/filepath"
	"io/fs"
	"strings"
)

type Markdown struct {
	Filepath string
	Content  string
}

type TagType string

const (
	Header1 TagType = "h1"
	Header2         = "h2"
	Header3         = "h3"
	Header4         = "h4"
	Paragraph       = "p"
)

type Tag struct {
	Type    TagType
	Content string
}

func read_markdowns(target_dir string) ([]*Markdown, error) {
	var out []*Markdown

	abs_target_dir, err := filepath.Abs(target_dir)
	if err != nil {
		return nil, err
	}

	err = filepath.WalkDir(abs_target_dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			bytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			out = append(out, &Markdown{
				Filepath: path,
				Content:  string(bytes),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}

func lex_line(line string) *Tag {
	words := strings.Fields(line)
	if len(words) < 2 {
		return &Tag{
			Type: Paragraph, 
			Content: line,
		}
	}

	head := words[0]
	content := strings.Join(words[1:], "")

	var tagType TagType
	if head == "#" {
		tagType = Header1 
	} else if head == "##" {
		tagType = Header2
	} else if head == "###" {
		tagType = Header3
	} else if head == "####" {
		tagType = Header4
	} else {
		tagType = Paragraph
	}

	return &Tag{
		Type: tagType, 
		Content: content,
	}
}

func lex_markdown(md *Markdown) []*Tag {
	var out []*Tag
	lines := strings.Split(md.Content, "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		ir := lex_line(line)
		out = append(out, ir)
	} 
	return out
}

func generate_tag_str(tag *Tag) string {
		out := "<" + string(tag.Type) + ">"
		out += tag.Content
		out += "</" + string(tag.Type) + ">"
		out += "\n"
		return out
}

func generate_html(tags []*Tag) string {
	i := 0
	html := ""

	for i < len(tags){
		tag := tags[i]
		switch (tag.Type) {
		case Header1, Header2, Header3, Header4: {
			html += generate_tag_str(tag)
		}
		case Paragraph: {
			html += generate_tag_str(tag)
		}
		}
		i += 1
	}
	return html
}

func main() {
	markdowns, err := read_markdowns("./testdata/") 
	if err != nil {
		log.Fatal(err)
	}

	for _, markdown := range markdowns {
		tags :=  lex_markdown(markdown)
		fmt.Println(generate_html(tags))
	}
}

