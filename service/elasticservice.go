

const (
	elasticUrl   = "http://localhost:9200"
	indexName    = "markhamca_idx"
	docType      = "web_page"
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