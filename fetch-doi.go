package main

import (
	"bufio"
	"fmt"
	"net/url" // Import the net/url package
	"os"
	"github.com/gocolly/colly"
)

func main() {
	// Create a new collector
	c := colly.NewCollector()

	// Take search input from the user
	fmt.Println("Enter your search query for Google Scholar:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() // Wait for input
	query := scanner.Text()

	// Encode the query to handle spaces and special characters
	encodedQuery := url.QueryEscape(query)

	// Construct the search URL for Google Scholar with the encoded query
	searchURL := fmt.Sprintf("https://scholar.google.com/scholar?q=%s", encodedQuery)

	// Setup Colly to parse the search results
	c.OnHTML("div.gs_ri", func(e *colly.HTMLElement) {
		title := e.ChildText("h3")
		authorAndPublicationInfo := e.ChildText(".gs_a")
		summary := e.ChildText(".gs_rs")
		fmt.Println("Title:", title)
		fmt.Println("Details:", authorAndPublicationInfo)
		fmt.Println("Summary:", summary)
		fmt.Println("------")
	})

	// Visit the search URL
	c.Visit(searchURL)
}
