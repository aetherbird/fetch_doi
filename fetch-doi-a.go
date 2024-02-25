package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"github.com/gocolly/colly"
	"github.com/manifoldco/promptui"
)

type Paper struct {
	Title string
	URL   string
}

func main() {
	// Initialize variables
	var papers []Paper

	// Create a new collector
	c := colly.NewCollector()

	// Setup Colly to parse the search results
	c.OnHTML("div.gs_ri", func(e *colly.HTMLElement) {
		title := e.ChildText("h3.gs_rt")
		link := e.ChildAttr("a", "href") // Assumes the first link in the div is the paper link

		papers = append(papers, Paper{Title: title, URL: link})
	})

	// Take search input from the user
	fmt.Println("Enter your search query for Google Scholar:")
	var query string
	fmt.Scanln(&query)

	// Encode the query
	encodedQuery := url.QueryEscape(query)

	// Construct the search URL
	searchURL := fmt.Sprintf("https://scholar.google.com/scholar?q=%s", encodedQuery)

	// Visit the search URL
	c.Visit(searchURL)

	// Ensure all HTTP requests have been completed
	c.Wait()

	// Create a slice of paper titles for the selectable list
	var titles []string
	for _, paper := range papers {
		titles = append(titles, paper.Title)
	}

	// Create a prompt with promptui
	prompt := promptui.Select{
		Label: "Select a paper",
		Items: titles,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	// Find the selected paper
	var selectedPaper Paper
	for _, paper := range papers {
		if paper.Title == result {
			selectedPaper = paper
			break
		}
	}

	// Fetch and display the selected paper's page
	resp, err := http.Get(selectedPaper.URL)
	if err != nil {
		fmt.Println("Error fetching paper:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Content of the selected paper's page:")
	fmt.Println(string(body))
}
