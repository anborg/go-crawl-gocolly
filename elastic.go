package main

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

func main() {

	// cert, _ := ioutil.ReadFile(*cacert)
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
			// "http://localhost:9201",
		},
		// Username: "foo",
		// Password: "bar",
		// CACert: cert,
		// Transport: &http.Transport{
		// 	MaxIdleConnsPerHost:   10,
		// 	ResponseHeaderTimeout: time.Second,
		// 	TLSClientConfig: &tls.Config{
		// 		MinVersion: tls.VersionTLS11,
		// 		// ...
		// 	},
		// 	// ...
		// }, //Transport
	}

	es, _ := elasticsearch.NewClient(cfg)

	log.Println(elasticsearch.Version)
	log.Println(es.Info())

}

func info(es *elasticapi) {
	
}
