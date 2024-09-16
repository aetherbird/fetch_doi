package main

import (
    "bufio"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "github.com/gocolly/colly"
    "github.com/rivo/tview"
)

type Paper struct {
    Title string
    URL   string
}

func main() {
    var papers []Paper  // list to store scraped paper details

    c := colly.NewCollector()  // create a new collector for web scraping

    // set up a callback to handle the elements matching the "div.gs_ri" selector
    c.OnHTML("div.gs_ri", func(e *colly.HTMLElement) {
        title := e.ChildText("h3.gs_rt")  // extract the paper's title
        link := e.ChildAttr("a", "href")  // extract the link to the paper
        papers = append(papers, Paper{Title: title, URL: link})  // add the paper to the list
    })

    fmt.Println("Enter your search query for Google Scholar:")  // prompt user for input
    scanner := bufio.NewScanner(os.Stdin)  // create a scanner to read user input
    scanner.Scan()  // scan for user input
    query := scanner.Text()  // store the input as a query

    encodedQuery := url.QueryEscape(query)  // encode the query for use in a URL
    searchURL := fmt.Sprintf("https://scholar.google.com/scholar?q=%s", encodedQuery)  // format the search url
    c.Visit(searchURL)  // visit the search URL to start scraping
    c.Wait()  // wait for the scraping process to complete

    app := tview.NewApplication()  // create a new terminal ui application
    list := tview.NewList()  // create a list widget to display papers

    for i, paper := range papers {  // loop through each scraped paper
        index := i  // capture the current value of i
        list.AddItem(paper.Title, "", 0, func() {  // add the paper title to the list
            selectedPaper := papers[index]  // retrieve the selected paper

            app.Stop()  // stop the ui application after selecting a paper

            // fetch and display the selected paper's web page
            resp, err := http.Get(selectedPaper.URL)  // send a GET request to the paper's URL
            if err != nil {  // check for errors during request
                fmt.Println("Error fetching paper:", err)  // print the error if request fails
                return
            }
            defer resp.Body.Close()  // ensure the response body is closed after use

            body, err := io.ReadAll(resp.Body)  // read the response body
            if err != nil {  // check for errors during reading
                fmt.Println("Error reading response body:", err)  // print the error if reading fails
                return
            }

            // print the content of the paper's page to the console
            fmt.Println("Content of the selected paper's page:")
            fmt.Println(string(body))
        })
    }

    // set up and run the ui application with the list widget
    if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
        panic(err)  // panic if there is an error running the ui application
    }
}
