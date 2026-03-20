package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Markdown struct {
	Filepath string
	Content  string
}

type TagType string

const (
	Header1   TagType = "h1"
	Header2   TagType = "h2"
	Header3   TagType = "h3"
	Header4   TagType = "h4"
	Paragraph TagType = "p"
	Image     TagType = "img"
)

type Tag struct {
	Type    TagType
	Content string
}

type ImageTag struct {
	Alt string
	Url string
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
			if filepath.Ext(path) != ".md" {
				return nil
			}
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
	if strings.HasPrefix(line, "!") {
		return &Tag{
			Type:    Image,
			Content: line[1:], // "[alt](url)"
		}
	}

	words := strings.Fields(line)
	if len(words) < 2 {
		return &Tag{
			Type:    Paragraph,
			Content: line,
		}
	}

	head := words[0]
	content := strings.Join(words[1:], " ")

	var tagType TagType
	switch head {
	case "#":
		tagType = Header1
	case "##":
		tagType = Header2
	case "###":
		tagType = Header3
	case "####":
		tagType = Header4
	default:
		tagType = Paragraph
	}

	return &Tag{
		Type:    tagType,
		Content: content,
	}
}

func lex_markdown(md *Markdown) []*Tag {
	var out []*Tag
	for line := range strings.SplitSeq(md.Content, "\n") {
		if len(line) == 0 {
			continue
		}
		ir := lex_line(line)
		out = append(out, ir)
	}
	return out
}

func parseImageTag(content string) (*ImageTag, error) {
	if !strings.HasPrefix(content, "[") {
		return nil, fmt.Errorf("Invalid image format")
	}
	content = strings.TrimPrefix(content, "[")

	alt, rest, ok := strings.Cut(content, "]")
	if !ok {
		return nil, fmt.Errorf("Invalid image format")
	}
	content = rest

	if !strings.HasPrefix(content, "(") {
		return nil, fmt.Errorf("Invalid image format")
	}
	content = strings.TrimPrefix(content, "(")

	url, _, ok := strings.Cut(content, ")")
	if !ok {
		return nil, fmt.Errorf("Invalid image format")
	}

	return &ImageTag{Alt: alt, Url: url}, nil
}

func generate_tag_str(tag *Tag) string {
	if tag.Type == Image {
		out := "<img "
		imgTag, err := parseImageTag(tag.Content)
		if err != nil {
			panic(err)
		}
		out += "src=\"" + imgTag.Url + "\" "
		out += "alt=\"" + imgTag.Alt + "\""
		out += ">"
		return out
	} else {
		out := "<" + string(tag.Type) + ">"
		out += tag.Content
		out += "</" + string(tag.Type) + ">"
		out += "\n"
		return out
	}
}

func generate_html(tags []*Tag) string {
	i := 0
	sb := strings.Builder{}

	for i < len(tags) {
		tag := tags[i]
		switch tag.Type {
		case Header1, Header2, Header3, Header4:
			{
				sb.WriteString(generate_tag_str(tag))
			}
		case Paragraph:
			{
				sb.WriteString(generate_tag_str(tag))
			}
		case Image:
			sb.WriteString(generate_tag_str(tag))
		}
		i += 1
	}
	return sb.String()
}

func main() {
	markdowns, err := read_markdowns("./testdata/")
	if err != nil {
		log.Fatal(err)
	}

	for _, markdown := range markdowns {
		tags := lex_markdown(markdown)
		fmt.Println(generate_html(tags))
	}
}
