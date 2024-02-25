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
    var papers []Paper

    c := colly.NewCollector()

    c.OnHTML("div.gs_ri", func(e *colly.HTMLElement) {
        title := e.ChildText("h3.gs_rt")
        link := e.ChildAttr("a", "href")
        papers = append(papers, Paper{Title: title, URL: link})
    })

    fmt.Println("Enter your search query for Google Scholar:")
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()
    query := scanner.Text()

    encodedQuery := url.QueryEscape(query)
    searchURL := fmt.Sprintf("https://scholar.google.com/scholar?q=%s", encodedQuery)
    c.Visit(searchURL)
    c.Wait()

    app := tview.NewApplication()
    list := tview.NewList()

    for i, paper := range papers {
        index := i // capture the current value of i
        list.AddItem(paper.Title, "", 0, func() {
            selectedPaper := papers[index]
            app.Stop()

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
        })
    }

    if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
        panic(err)
    }
}
