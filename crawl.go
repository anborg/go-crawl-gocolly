package main

import (
	"encoding/json"
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

	baseUrl := "https://www.coursera.org"

	// q, _ := queue.New(
	// 	2, // Number of consumer threads
	// 	&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	// )

	c := colly.NewCollector(
		// Visit only domains: coursera.org, www.coursera.org
		colly.AllowedDomains("www.coursera.org"),

		//colly.AllowURLRevisit(),

		// MaxDepth is 1, so only the links on the scraped page
		// is visited, and no further links are followed
		//colly.MaxDepth(1),

		//#FORPARALLEL-code-anb
		colly.Async(true),

		// Attach a debugger to the collector
		// colly.Debugger(&debug.LogDebugger{}), // "github.com/gocolly/colly/debug"
	)
	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	//#FORPARALLEL-code-anb  http://go-colly.org/docs/examples/parallel/
	// Limit the number of threads to 2
	// when visiting links which domains' matches "*httpbin.*" glob
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		//  RandomDelay: 5 * time.Second,
		//Delay:      5 * time.Second,
	})

	// authenticate if required
	// err := c.Post("http://example.com/login", map[string]string{"username": "admin", "password": "admin"})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// log.Println(fmt.Sprintf("\tTO CRAWL  : %s", link))
		//Skip the following
		//#, javascript:; ,  #

		// start scaping the page under the link found
		e.Request.Visit(link)

		// Add URLs to the queue
		// q.AddURL(link)

		//c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnHTML("title", func(e *colly.HTMLElement) {
		//fmt.Println(e.Text)
		//e.Ctx.Put("title", e.Text)
	})

	pageCount := 0
	c.OnRequest(func(r *colly.Request) {

		// log.Println(fmt.Sprintf("%d  Visiting : %s", pageCount, r.URL))

		// Before making a request put the URL with
		// the key of "url" into the context of the request
		r.Ctx.Put("url", r.URL.String())
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Print the response
	c.OnResponse(func(r *colly.Response) {
		pageCount++
		//	log.Printf("%s\n", bytes.Replace(r.Body, []byte("\n"), nil, -1))
		// log.Println(string(r.Body))
		//log.Println(r.Headers)

		// After making a request get "url" from
		// the context of the request
		urlVisited := r.Ctx.Get("url")
		webpage := WebPage{
			Url:     r.Ctx.Get("url"),      //- can be put in ctx c.OnRequest, and r.Ctx.Get("url")
			Title:   "my page title",       //string(r.title), // Where to get this?
			Content: "page html body bla ", //string(r.Body),   //string(r.Body) - can be done in c.OnResponse
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		// Dump json to the standard output
		enc.Encode(webpage)

		log.Println(fmt.Sprintf("%d  DONE Visiting : %s", pageCount, urlVisited))

	})

	c.Visit(baseUrl)
	//#FORPARALLEL-code-anb
	// Wait until threads are finished
	c.Wait()

	// q.AddURL(baseUrl)
	// q.Run(c)
}
