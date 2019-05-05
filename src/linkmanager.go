package main

import (
	"net/url"
	"regexp"
	"strings"
)

//DivideURLByType divide link by external and internal
func DivideURLByType(links []string, domainUrl string) (map[string]int, map[string]int) {
	var external = make(map[string]int)
	var internal = make(map[string]int)

	isValid, _ := regexp.Compile("^(http[s]?:\\/\\/)(((www)?)([a-z0-9A-Z\\-\\.]+))\\/([a-zA-Z0-9\\/\\-\\_\\:\\+\\.]+)?$")
	isInternal, _ := regexp.Compile("^(http[s]?:\\/\\/)(www.)?" + domainUrl + "\\/(.*)$")

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

/**
* Get Cleaned Url from given url
* @parameter
* @return
 */
func GetCleanUrl(urlString string) string {
	u, err := url.Parse(urlString)
	checkErr(err)
	return u.Scheme + "://" + u.Host
	// regex, err := regexp.Compile("^(http[s]?:\\/\\/)(((www)?)([a-z0-9A-Z\\-\\.]+))((\\/)[a-zA-Z0-9\\/\\-\\_\\:\\+\\.]+)?$")
	// checkErr(err)
	// found := regex.FindStringSubmatch(url)
	// return found[1] + found[2]
}

/**
* Get domain from given url
* @parameter
* @return
 */
func GetDomainFromUrl(urlString string) string {
	u, err := url.Parse(urlString)
	checkErr(err)
	return u.Host
	// regex, err := regexp.Compile("(http(s)?:\\/\\/)([\\w.-]+)(\\.[\\w\\.-]+)")
	// checkErr(err)
	// found := regex.FindStringSubmatch(startUrl)
	// return found[3] + found[4]
}

/**
* Get map of clean url  from given map of url
* @parameter
* @return
 */
func GetCleanedURLMap(urls map[string]int) map[string]int {
	var returnMap = make(map[string]int)
	for key := range urls {
		if _, ok := returnMap[key]; !ok {
			returnMap[GetCleanUrl(key)] = 1
		}
	}
	return returnMap
}
