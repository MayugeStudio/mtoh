package main

import (
	"fmt"
	"os"
	"log"
	"path/filepath"
	"io/fs"
)

type Markdown struct {
	Filepath string
	Content  string
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

func main() {
	markdowns, err := read_markdowns("./testdata/") 
	if err != nil {
		log.Fatal(err)
	}

	for _, markdown := range markdowns {
		fmt.Println(markdown.Filepath)
		fmt.Println(markdown.Content)
	}
}

