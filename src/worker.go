package main

import (
	"./DomManager"
	"./HttpRequester"
	"./databases"
	"./memqueue"
	"os"
	"strconv"
	"sync"
	"time"
)

func worker(jobs <-chan string, dbManager databases.DataManager, queue *memqueue.Queue, recent *memqueue.Queue, w *widgets) {
	for next := range jobs {
		if next == "" {
			continue
		}
		printStart(next, w)
		/*Prepare request*/
		requester := HttpRequester.HttpRequester{
			URL:       next,
			UserAgent: "Url Scraper AC",
			Timeout:   4,
		}
		/*Actual Check*/
		document, err := requester.MakeCheck()
		if checkRequestError(err, dbManager, next, recent, w) {
			continue
		}

		/*Get Link from response document*/
		dom := DomManager.DomManager{
			Document: document,
		}
		urls := dom.GetURLFromDocument()
		extMap, intMap := DivideURLByType(urls, GetDomainFromURL(next))
		storeResponse(next, extMap, intMap, dbManager, queue, recent)

		printComplete(w)
		wait, _ := strconv.Atoi(os.Getenv("WAITING_TIME"))
		time.Sleep(time.Second * time.Duration(wait))
	}

}

func storeResponse(current string, extMap map[string]int, intMap map[string]int, dbManager databases.DataManager, queue *memqueue.Queue, recent *memqueue.Queue) {
	var wg sync.WaitGroup

	/*Convert Link From Map To Slice (Array)*/
	var external []string
	var internal []string

	wg.Add(2)
	go func() {
		external = mapToSlice(extMap) //Get a slice from a map
		wg.Done()
	}()
	go func() {
		internal = mapToSlice(intMap) //Get a slice from a map
		wg.Done()
	}()
	_ = internal //Remove Declared but not used error

	go func() {
		//Insert DB to mysql
		dbManager.RegisterURLToMysql(GetCleanedURLMap(extMap))
	}()

	wg.Wait() //Wait mapToSlice to be completed

	wg.Add(3)
	go func() {
		/*Save to redis*/
		dbManager.RegisterURLToRedis(external)
		wg.Done()
	}()
	go func() {
		/*Save in Queue*/
		for _, v := range external {
			queue.Push(memqueue.URL{URL: v})
		}
		wg.Done()
	}()
	go func() {
		/*Add done URL to last URL Memory Queue*/
		recent.Push(memqueue.URL{
			URL:    current,
			Domain: GetDomainFromURL(current),
			Clean:  GetCleanURL(current),
		})
		wg.Done()
	}()
	wg.Wait()
}
