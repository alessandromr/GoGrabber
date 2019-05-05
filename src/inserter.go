package main

import (
	"./databases"
	"./memqueue"
	"log"
	"time"
)

func inserter(jobs chan<- string, dbManager databases.DataManager, queue *memqueue.Queue, recent *memqueue.Queue) {
	for {
		time.Sleep(time.Millisecond * 500)

		//If there are no recent elements just get one from queue
		if len(recent.URLs) < 1 {
			mess, err := queue.Pop()
			if err == nil {
				jobs <- mess.URL
			}
			continue
		}

		//Peek last n elements from queue
		messages, err := queue.PeekN(recent.Qlen)
		if err != nil {
			log.Println("No elements in queue")
			continue
		}

		var good string
		bad := false

		//Check if one element in queue is not a recent one
		for _, inQueue := range messages {
			bad = false
			for _, recent := range recent.URLs {
				if GetDomainFromURL(inQueue.URL) == recent.Domain {
					bad = true
				}
			}
			if !bad {
				good = inQueue.URL
				queue.Delete(good)
				jobs <- good
				continue
			}
		}
		//If no element is found return empty
		if good == "" {
			queue.Trim(recent.Qlen)
			log.Println("No new URLs in queue, trimming it and trying again")
			continue
		}
	}
}
