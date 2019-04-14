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
		log.Println("Checking for new jobs")
		//If there are no recent elements just get one from queue
		if len(recent.URLs) < 1 {
			mess, _ := queue.Pop()
			jobs <- mess.URL
			continue
		}

		//Peek last n elements from queue
		messages, _ := queue.PeekN(recent.Qlen)

		var good string
		bad := false

		//Check if one element in queue is not a recent one
		for _, inQueue := range messages {
			bad = false
			for _, recent := range recent.URLs {
				if GetDomainFromUrl(inQueue.URL) == recent.Domain {
					bad = true
				}
			}
			if !bad {
				good = inQueue.URL
				queue.Delete(good)
				log.Println("Insert: ", good)
				jobs <- good
				break
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
