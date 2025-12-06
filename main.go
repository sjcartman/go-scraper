package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

func main() {
	resp, err := http.Get("https://books.toscrape.com/")
	if err != nil {
		log.Fatalf("Request failed")
	}
	defer resp.Body.Close()
	if err != nil {
		log.Fatal("Error Reading Response Body")
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal("Error Decoding Body")
	}
	exploreNode(doc)

}

func exploreNode(node *html.Node) {
	for _, item := range node.Attr {
		if item.Key == "href" {
			fmt.Println(item.Val)
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		exploreNode(child)
	}
}
