package main

import (
	"fmt"
	"os"
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
	UnorderedList   = "ul"
	OrderedList     = "ol"
	ListItem        = "li"
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
	} else if head == "-" {
		tagType = ListItem
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

func main() {
	markdowns, err := read_markdowns("./testdata/") 
	if err != nil {
		log.Fatal(err)
	}

	for _, markdown := range markdowns {
		for _, tag := range lex_markdown(markdown) {
			fmt.Println(tag.Type, tag.Content)
		}
	}
}

