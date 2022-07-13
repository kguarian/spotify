package main

import (
	"flag"
	"io/fs"
	"log"
	"os"
)

var (
	replace   = flag.Bool("replace", false, "")
	tgt_dir   = flag.String("tgt", "", "")
	pre_text  = flag.String("pre", "", "")
	post_text = flag.String("post", "", "")
)

func main() {
	flag.Parse()

	if *replace {
		//inductive problem. replace nothing with something? Infinite file?
		if *pre_text == "" || *tgt_dir == "" {
			flag.Usage()
			log.Fatalf("replace requires pre-text and target directory")
		}
		//replace in directory
		dir, err := os.Open(*tgt_dir)
		if err != nil {
			log.Fatalf("error opening target directory: %v", err)
		}
		defer dir.Close()
		files, err := dir.Readdir(0)
		if err != nil {
			log.Fatalf("error reading target directory: %v", err)
		}
		replaceInDirectory(files, *pre_text, *post_text)

	}
}

func replaceInDirectory(files []fs.FileInfo, preText, postText string) {
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		f, err := os.Open(file.Name())
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		dirContents, err := f.Readdir(0)
		if err != nil {
			log.Fatalf("error reading target directory: %v", err)
		}
		replaceInDirectory(dirContents, preText, postText)
	}
}
