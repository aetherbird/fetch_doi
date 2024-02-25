package main

import (
    "fmt"
    "net/http"

    "github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://aethercosmology.com/")

  c.OnHTML("body", func(e *colly.HTMLElement) {
      // Extract text using various methods:
      text := e.Text()           // Get all text content
      title := e.ChildText("h1") // Get text of H1 element
      paragraphs := e.ChildTexts("p") // Get text of all P elements
  
      // Print the extracted text:
      fmt.Println("Title:", title)
      fmt.Println("Paragraphs:")
      for _, paragraph := range paragraphs {
          fmt.Println("- " + paragraph)
      }
  })

}
