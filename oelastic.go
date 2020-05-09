package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
	"github.com/olivere/elastic/v7"
)

func main(){
	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200")
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetRetrier(NewCustomRetrier()),
		elastic.SetGzip(true),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
		elastic.SetHeaders(http.Header{
		  "X-Caller-Id": []string{"..."},
		}),
	)
	if err != nil {
	// Handle error
	panic(err)
	}
	defer client.Stop()

	exists, err := client.IndexExists("websiteindex").Do(context.Background())
	if err != nil {
		// Handle error
	}
	if exists {
		println("Index exists", "websiteindex")
	}


	// Create a new index.
mapping := `{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"properties":{
			"url"		:{"type":"string", "index" : "not_analyzed" },
			"title"		:{"type":"string"},
			"content"	:{"type":"string"},
			"time"		:{"type":"date"},
		}
	}
}`

ctx := context.Background()
createIndex, err := client.CreateIndex("websiteindex").BodyString(mapping).Do(ctx)
if err != nil {
    // Handle error
    panic(err)
}
if !createIndex.Acknowledged {
    // Not acknowledged
}

}