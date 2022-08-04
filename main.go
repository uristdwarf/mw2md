package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var outputDir = "output"

type MWExport struct {
	XMLName xml.Name   `xml:"mediawiki"`
	Page    []WikiText `xml:"page"`
}

type WikiText struct {
	XMLName xml.Name `xml:"page"`
	Title   string   `xml:"title"`
	Text    []byte   `xml:"revision>text"`
}

type Text struct {
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("USAGE: wiki2md [INPUT]")
		os.Exit(1)
	}
	input, err := os.ReadFile(os.Args[1])
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("read error: File not found: %v", os.Args[1])
		os.Exit(2)
	} else if err != nil {
		log.Fatal(err)
	}

	err = os.Mkdir(outputDir, 0777)
	if !errors.Is(err, os.ErrExist) && err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	var mwBackup MWExport
	if err = xml.Unmarshal(input, &mwBackup); err != nil {
		log.Fatal(err)
	}

	for _, content := range mwBackup.Page {
		path := ""
		if strings.Contains(content.Title, "/") {
			content.Title = getPath(content.Title)
		}
		log.Printf("Converting %v\n", content.Title)
		md, err := convert(content.Text)
		if err != nil {
			log.Printf("failed to convert: %v", err)
			continue
		}
		createFile(content.Title, outputDir+"/"+path, md)
	}
}

func getPath(title string) string {
	return strings.ReplaceAll(title, "/", "-")

}

func convert(content []byte) ([]byte, error) {
	cmd := exec.Command("pandoc", "-f", "mediawiki", "-t", "markdown")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return []byte(""), err
	}

	go func() {
		defer stdin.Close()
		stdin.Write(content)
	}()

	return cmd.Output()
}

func createFile(title string, path string, content []byte) {
	if path != "" {
		if err := os.MkdirAll(filepath.Dir(path), 0660); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	if err := os.WriteFile(path+title+".md", []byte(content), 0660); err != nil {
		fmt.Println(err)
		return
	}
}
