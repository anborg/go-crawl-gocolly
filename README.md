# crawl-gocolly

TODO:
- create a flag to do indexing explicitly e.g -index http://localhost:9200
- neatly seaprate elastic from crawling




go build crawl.go
./crawl -url https://www.markham.ca/wps/portal/home

go run .\crawl.go -url https://www.markham.ca/wps/portal/home


## NEXT:
- searchui
http://localhost:9200/markhamca_idx/_search?q=Markham&pretty

http://localhost:9200/_cat/indices?v



## DONE
- crawl http
- index ES - with oliver-elastic


## TODO 
- setup search UI
- cleanup js/html
- parse pdf, word, xls

- ES-index with https://github.com/elastic/go-elasticsearch


Read
https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl.html
https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-function-score-query.html

#Go cmd line
https://github.com/spf13/cobra - CLI interface

