package main

import (
	"./memqueue"
	"github.com/mum4k/termdash/terminal/termbox"
	"log"
	"time"
)

func supervisor(recent *memqueue.Queue, queue *memqueue.Queue, t *termbox.Terminal) {
	for {
		time.Sleep(time.Millisecond * 600)
		if !checkRecent(recent) {
			t.Close()
			log.Fatalln("Too many duplicates")
		}
		if !checkQueue(queue) {
			t.Close()
			log.Fatalln("Too many duplicates")
		}
	}
}

/*
	Remove duplicates from queue
*/
func checkQueue(queue *memqueue.Queue) bool {
	//TODO
	return true
}

/*
	Check recent for duplicates
*/
func checkRecent(recent *memqueue.Queue) bool {
	var duplicates int
	for key, toCheck := range recent.URLs {
		tmp := toCheck.Domain
		for targetKey, target := range recent.URLs {
			if target.Domain == tmp {
				if targetKey != key {
					duplicates++
				}
			}
		}
	}
	var max int
	if len(recent.URLs) < 15 {
		max = 3 //If duplicates are more than 3
	}
	if len(recent.URLs) > 15 {
		max = int(float64(len(recent.URLs)) * 0.2) //If duplicates are more than the 20% of the total queue length
	}
	if duplicates > max {
		// t.Close()
		// log.Println(recent.URLs)
		// log.Println("max", max)
		// log.Println("dup", duplicates)
		// log.Fatal()
		return false
	}
	return true
}