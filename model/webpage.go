package domain

import (
	"fmt"
	"time"
)

type WebPage struct {
	Url     string    `json:"url"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
}

func (document *WebPage) Print() {
	// enc := json.NewEncoder(os.Stdout)
	// enc.SetIndent("", "  ")
	// document.Content = ""
	// enc.Encode(document)
	println(fmt.Sprintf("page:  {\n  title: %s, \n  url : %s, \n  content:%s, \n  time:%s \n}", document.Title, document.Url, "-redacted-", document.Time))
}
