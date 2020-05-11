package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"crawl/domain"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/olivere/elastic/v7"
)

const (
	elasticUrl   = "http://localhost:9200"
	indexName    = "markhamca_idx"
	// docType      = "webpages"
	indexMapping = `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 1
		},
		// "markhamca_idx" : {
			"mappings" : {
				"properties" : {
					"url" : { "type" : "text" 	},
					"title" : { "type" : "text"	},
					"content" : { "type" : "text"  },
					"time" : { "type" : "date" }
				}
			}
		//}
	}`
)

func main() {
	var baseUrl string

	flag.StringVar(&baseUrl, "url", "", "Crawl url is required")
	flag.Parse()

	if baseUrl == "" {
		log.Println("crawl url required")
		os.Exit(1)
	}

	//==================
	// var client elastic.Client
	client, err := elastic.NewClient(
		elastic.SetURL(elasticUrl),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		// elastic.SetRetrier(NewCustomRetrier()),
		elastic.SetGzip(true),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
		elastic.SetHeaders(http.Header{
			"X-Caller-Id": []string{"..."},
		}),
	)
	if err != nil {
		panic(err)
	}
	// _ = info(client)
	deleteIndex(client)
	err = createIndex(client)
	if err != nil {
		panic(err)
	}

	document := domain.WebPage{}
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
		document.Url = r.URL.String()

	})

	c.OnResponse(func(r *colly.Response) { //get body
		pageCount++
		// document.Content = "page html body bla " // string(r.Body)
		//urlVisited := r.Ctx.Get("url")
		println(fmt.Sprintf("  DONE Visiting %d: %s", pageCount, document.Url))
		document.Print()

	})

	c.OnHTML("html head title", func(e *colly.HTMLElement) { // Title
		//e.Ctx.Put("title", e.Text)
		document.Title = e.Text
	})
	c.OnHTML("html body", func(e *colly.HTMLElement) { // Body / content
		//e.Ctx.Put("title", e.Text)
		document.Content = e.Text
		document.Time = time.Now()
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
		// enc := json.NewEncoder(os.Stdout)
		// enc.SetIndent("", "  ")
		// enc.Encode(document)
		insertDocument(client, document)
		document.Print()
	})

	//#FORPARALLEL-code-anb
	// Wait until threads are finished
	c.Visit(baseUrl)
	c.Wait()

	// q.AddURL(baseUrl)
	// q.Run(c)
}

func insertDocument(client *elastic.Client, document domain.WebPage) error {
	_, err := client.Index().
		Index(indexName).
		// Type(docType).
		BodyJson(document).
		// Id(document.Url). //Make id automatic - dupes check
		Do(context.Background())

	if err != nil {
		return err
	}
	_, _ = client.Flush().Index(indexName).Do(context.Background())
	return nil
}

func createIndex(client *elastic.Client) error {
	exists, err := client.IndexExists(indexName).Do(context.Background())
	if err != nil {
		return err
	}
	if exists {
		fmt.Printf("Index already exists : %s\n", indexName)
		return nil
	}

	res, err := client.CreateIndex(indexName).
		Body(indexMapping).
		// Type(webpages).
		Do(context.Background())

	if err != nil {
		return err
	}
	if !res.Acknowledged {
		return errors.New("CreateIndex was not acknowledged. Check that timeout value is correct.")
	}
	fmt.Printf("Index newly Created  : %s\n", indexName)

	// return addWebPagesToIndex(client)
	return nil
}
func deleteIndex(client *elastic.Client) error {
	_, _ = client.DeleteIndex(indexName).Do(context.Background())
	return nil
} //deleteIndex
