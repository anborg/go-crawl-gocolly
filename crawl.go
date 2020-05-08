package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

type WebPage struct {
	Url     string `json:"url"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func main() {
	var baseUrl string

	flag.StringVar(&baseUrl, "url", "", "Crawl url is required")
	flag.Parse()

	if baseUrl == "" {
		log.Println("crawl url required")
		os.Exit(1)
	}

	webpage := WebPage{}
	// baseUrl := "https://www.coursera.org"
	// baseUrl := "https://www.markham.ca/wps/portal/home"

	pageCount := 0

	c := colly.NewCollector(
		colly.AllowedDomains("www.markham.ca"),
		//colly.AllowURLRevisit(),
		//colly.MaxDepth(1),
		colly.Async(true),
		// colly.Debugger(&debug.LogDebugger{}), // "github.com/gocolly/colly/debug"
		//colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)
	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2, // http://go-colly.org/docs/examples/parallel/
		//  RandomDelay: 5 * time.Second,
		//Delay:      5 * time.Second,
	})

	// q, _ := queue.New(
	// 	2, // Number of consumer threads
	// 	&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	// )
	// q.AddURL(baseUrl)

	// authenticate if required
	// err := c.Post("http://example.com/login", map[string]string{"username": "admin", "password": "admin"})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	c.OnRequest(func(r *colly.Request) { // url
		webpage.Url = r.URL.String()

	})

	c.OnResponse(func(r *colly.Response) { //get body
		pageCount++
		//webpage.Content = "page html body bla " // string(r.Body)
		//urlVisited := r.Ctx.Get("url")
		println(fmt.Sprintf("  DONE Visiting %d: %s", pageCount, webpage.Url))

	})

	c.OnHTML("html head title", func(e *colly.HTMLElement) { // Title
		//e.Ctx.Put("title", e.Text)
		webpage.Title = e.Text
	})
	c.OnHTML("html body", func(e *colly.HTMLElement) { // Body / content
		//e.Ctx.Put("title", e.Text)
		webpage.Content = "my body bla" // e.Text
	})

	// On every a element which has href attribute
	c.OnHTML("a[href]", func(e *colly.HTMLElement) { // href , callback
		link := e.Attr("href")
		e.Request.Visit(link)
		// q.AddURL(e.Request.AbsoluteURL(link)) // Add URLs to the queue
		//c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnError(func(r *colly.Response, err error) { // Set error handler
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnScraped(func(r *colly.Response) { // DONE
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(webpage)
	})

	//#FORPARALLEL-code-anb
	// Wait until threads are finished
	c.Visit(baseUrl)
	c.Wait()

	// q.AddURL(baseUrl)
	// q.Run(c)
}
