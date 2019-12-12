package  dommanager

import (
	"github.com/PuerkitoBio/goquery"
)

//DomManager manage DOM JQuery style
type DomManager struct {
	Document goquery.Document
}

//GetURLFromDocument Get Links from a given document
func (dm DomManager) GetURLFromDocument() []string {
	var links []string
	dm.Document.Find("a").Each(func(index int, element *goquery.Selection) {
		href, exists := element.Attr("href")
		if exists {
			links = append(links, href)
		}
	})
	return links
}
