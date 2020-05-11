package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	// "gopkg.in/olivere/elastic.v7"
	"github.com/olivere/elastic/v7"
)

const (
	elasticUrl   = "http://localhost:9200"
	indexName    = "markhamca_idx"
	docType      = "web_page"
	queryString  = "My WebPage content 7"
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

type WebPage struct {
	Url     string    `json:"url"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
}

func main() {
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

	_ = info(client)
	// deleteIndex(client)

	// err = createIndex(client)
	// if err != nil {
	// 	panic(err)
	// }

	// err = insertDocument(client)
	// if err != nil {
	// 	panic(err)
	// }
	err = matchQueryExample(client)
	if err != nil {
		panic(err)
	}
	// err = idQueryExample(client)
	// if err != nil {
	// 	panic(err)
	// }
	// err = termQueryExample(client)
	// if err != nil {
	// 	panic(err)
	// }
} // main()

func info(client *elastic.Client) error {
	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping(elasticUrl).Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion(elasticUrl)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)
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

func insertDocument(client *elastic.Client) error {
	for i := 0; i < 10; i++ {
		document := WebPage{
			Url:     fmt.Sprintf("http://mywebsite.com/%d", i),
			Title:   fmt.Sprintf("My WebPage %d", i),
			Content: fmt.Sprintf("My WebPage content %d", i),
			Time:    time.Now(),
		}

		_, err := client.Index().
			Index(indexName).
			// Type(docType).
			BodyJson(document).
			Do(context.Background())

		if err != nil {
			return err
		}

	} //for

	// Flush to make sure the documents got written.
	_, _ = client.Flush().Index(indexName).Do(context.Background())

	return nil
}

func matchQueryExample(esclient *elastic.Client) error {

	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchQuery("title", "My WebPage 7"))
	printQuery(searchSource)

	searchService := esclient.Search().Index(indexName).SearchSource(searchSource)
	searchResult, err := searchService.Do(context.Background())
	if err != nil {
		fmt.Println("[ProductsES][GetPIds]Error=", err)
		return err
	}
	// printSearchResult(searchResult)
	var webpageArray []WebPage
	for _, hit := range searchResult.Hits.Hits {
		var webpage WebPage
		err := json.Unmarshal(hit.Source, &webpage)
		if err != nil {
			fmt.Println("[Getting webpage][Unmarshal] Err=", err)
		}
		webpageArray = append(webpageArray, webpage)
	}

	if err != nil {
		fmt.Println("Fetching student fail: ", err)
		return err
	} else {
		for _, page := range webpageArray {
			fmt.Printf("Webpage found url: %s, title: %d, content: %f \n", page.Url, page.Title, page.Content)
		}
	}
	return nil
}

func printQuery(searchSource *elastic.SearchSource) {
	/* this block will basically print out the es query */
	queryStr, err1 := searchSource.Source()
	queryJs, err2 := json.Marshal(queryStr)

	if err1 != nil || err2 != nil {
		fmt.Println("[esclient][GetResponse]err during query marshal=", err1, err2)

	}
	fmt.Println("[esclient]Final ESQuery=\n", string(queryJs))

}

// func printSearchResult(result *client.SearchResult) {

// var webpageArray []WebPage
// for _, hit := range searchResult.Hits.Hits {
// 	var webpage WebPage
// 	err := json.Unmarshal(hit.Source, &webpage)
// 	if err != nil {
// 		fmt.Println("[Getting webpage][Unmarshal] Err=", err)
// 	}
// 	webpageArray = append(webpageArray, webpage)
// }

// if err != nil {
// 	fmt.Println("Fetching student fail: ", err)
// } else {
// 	for _, page := range webpageArray {
// 		fmt.Printf("Webpage found url: %s, title: %d, content: %f \n", page.Url, page.Title, page.Content)
// 	}
// }
// 	return nil
// }

func idQueryExample(client *elastic.Client) error {
	docId := "c9f367d58ade70c0f68f2ad0382fcfd2998a32b0"
	// Get tweet with specified ID
	get1, err := client.Get().
		Index(indexName).
		Type(docType).
		Id(docId).
		Do(context.Background())
	if err != nil {
		switch {
		case elastic.IsNotFound(err):
			panic(fmt.Sprintf("Document not found: %v", err))
		case elastic.IsTimeout(err):
			panic(fmt.Sprintf("Timeout retrieving document: %v", err))
		case elastic.IsConnErr(err):
			panic(fmt.Sprintf("Connection problem: %v", err))
		default:
			// Some other kind of error
			panic(err)
		}
	}
	fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)

	return nil
} //termQueryExample

func deleteIndex(client *elastic.Client) error {
	_, _ = client.DeleteIndex(indexName).Do(context.Background())
	return nil
} //deleteIndex

func termQueryExample(client *elastic.Client) error {
	// Search with a term query
	termQuery := elastic.NewTermQuery("content", queryString)
	searchResult, err := client.Search().
		Index(indexName).        // search in index "twitter"
		Query(termQuery).        // specify the query
		Sort("url", true).       // sort by "user" field, ascending
		From(0).Size(10).        // take documents 0-9
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		// Handle error
		panic(err)
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over iterating the hits, see below.
	var ttyp WebPage
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		t := item.(WebPage)
		fmt.Printf("Tweet by %s: %s\n", t.Url, t.Content)
	}
	// TotalHits is another convenience function that works even when something goes wrong.
	fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())

	// Here's how you iterate through results with full control over each step.
	// if searchResult.Hits.TotalHits > 0 {
	// 	fmt.Printf("Found a total of %d tweets\n", searchResult.Hits.TotalHits)

	// 	// Iterate through results
	// 	for _, hit := range searchResult.Hits.Hits {
	// 		// hit.Index contains the name of the index

	// 		// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
	// 		var t WebPage
	// 		err := json.Unmarshal(*hit.Source, &t)
	// 		if err != nil {
	// 			// Deserialization failed
	// 		}

	// 		// Work with tweet
	// 		fmt.Printf("Tweet by %s: %s\n", t.Url, t.Content)
	// 	}
	// } else {
	// 	// No hits
	// 	fmt.Print("Found no tweets\n")
	// }

	return nil
} //queryWebPages
