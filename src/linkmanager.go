package main

import (
	"net/url"
	"regexp"
	"strings"
)

//DivideURLByType divide link by external and internal
func DivideURLByType(links []string, domainURL string) (map[string]int, map[string]int) {
	var external = make(map[string]int)
	var internal = make(map[string]int)

	isValid, _ := regexp.Compile("^(http[s]?:\\/\\/)(((www)?)([a-z0-9A-Z\\-\\.]+))\\/([a-zA-Z0-9\\/\\-\\_\\:\\+\\.]+)?$")
	isInternal, _ := regexp.Compile("^(http[s]?:\\/\\/)(www.)?" + domainURL + "\\/(.*)$")

	for _, element := range links {
		element = strings.TrimSuffix(element, "/")
		if isValid.MatchString(element) {
			if isInternal.MatchString(element) {
				if _, ok := internal[element]; !ok {
					internal[element] = 1
				}
			} else {
				if _, ok := external[element]; !ok {
					external[element] = 1
				}
			}
		}
	}
	return external, internal
}

//GetCleanURL return cleaned Url from given url
func GetCleanURL(urlString string) string {
	u, err := url.Parse(urlString)
	checkErr(err)
	return u.Scheme + "://" + u.Host
}

//GetDomainFromURL return domain from given url
func GetDomainFromURL(urlString string) string {
	u, err := url.Parse(urlString)
	checkErr(err)
	return u.Host
}

//GetCleanedURLMap return map of clean URLs from given map of URLs
func GetCleanedURLMap(urls map[string]int) map[string]int {
	var returnMap = make(map[string]int)
	for key := range urls {
		if _, ok := returnMap[key]; !ok {
			returnMap[GetCleanURL(key)] = 1
		}
	}
	return returnMap
}
