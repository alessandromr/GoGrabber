package main

import (
	"./databases"
	"./memqueue"
	"log"
	"os"
	"time"
)

func main() {
	loadEnv()
	//Prepare the Database clients
	var dbManager = databases.DataManager{}
	dbManager.SetClients()

	queue := memqueue.Queue{Qlen: 1500}
	recent := memqueue.Queue{Qlen: 40}
	jobs := make(chan string, 100)

	jobs <- os.Getenv("START_URL")

	go worker(jobs, dbManager, &queue, &recent)
	go inserter(jobs, dbManager, &queue, &recent)

	go func() {
		for {
			log.Println("")
			log.Println("Status:")
			log.Println("Queue len: ", len(queue.URLs))
			log.Println("Pending jobs: ", len(jobs))
			log.Println("")
			time.Sleep(time.Second * 5)
		}
	}()

	for {

	}

}
