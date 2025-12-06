package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"golang.org/x/net/html"
)

type Book struct {
	price float64
}

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
	books := []Book{}
	exploreBooks(doc, &books)
	fmt.Println(books)

}

func exploreBooks(node *html.Node, books *[]Book) {
	if searchAttributes(node, "class", "col-xs-6 col-sm-4 col-md-3 col-lg-3") {
		book := Book{}
		parseBookNode(node, &book)
		*books = append(*books, book)
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		exploreBooks(child, books)
	}
}

func searchAttributes(node *html.Node, key string, val string) bool {
	for _, item := range node.Attr {
		if item.Key == key && item.Val == val {
			return true
		}
	}
	return false
}

func parseBookNode(node *html.Node, book *Book) {
	if searchAttributes(node, "class", "price_color") {
		price, err := parsePriceNode(node)
		if err != nil {
			fmt.Println(err)
		} else if book.price > 0.0 {
			fmt.Printf("error Book all ready has a price")
		} else {
			book.price = price
		}

	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		parseBookNode(child, book)
	}

}

func parsePriceNode(node *html.Node) (float64, error) {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		price, err := parsePrice(child.Data)
		if err != nil {
			log.Println(err)
			continue
		}
		return price, nil

	}
	return 0, fmt.Errorf("no price found in node: %v", node.Attr)

}

func parsePrice(raw string) (float64, error) {
	re := regexp.MustCompile(`\d+\.\d+`)
	priceStr := re.FindString(raw)
	if priceStr == "" {
		return 0, fmt.Errorf("no price found in string: %s", raw)
	}
	return strconv.ParseFloat(priceStr, 64)
}
