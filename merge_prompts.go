package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	dirs := []string{"general", "libs"}
	var files []string

	for _, dir := range dirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".md" {
				files = append(files, path)
			}
			return nil
		})
	}

	sort.Strings(files)

	outFile, err := os.Create("prompt.md")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Printf("Failed to read file %s: %v", file, err)
			continue
		}
		_, err = outFile.Write(content)
		if err != nil {
			log.Printf("Failed to write to prompt.md: %v", err)
			continue
		}
		// Add a newline between files for separation
		_, err = outFile.WriteString("\n")
		if err != nil {
			log.Printf("Failed to write to prompt.md: %v", err)
		}
	}

	fmt.Println("Successfully merged markdown files to prompt.md")
}
