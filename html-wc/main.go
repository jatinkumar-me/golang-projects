package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	words, images := countWordsAndImages(doc)
	fmt.Printf("%d words and %d images\n", words, images)
}

func countWordsAndImages(node *html.Node) (int, int) {
	var words, images int
	visit(&words, &images, node)
	return words, images
}

func visit(words, images *int, node *html.Node) {
	if node.Type == html.TextNode {
		*words += len(strings.Fields(node.Data))
	} else if node.Type == html.ElementNode && node.Data == "img" {
		*images++
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		visit(words, images, child)
	}
}
