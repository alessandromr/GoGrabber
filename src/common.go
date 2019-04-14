package main

import (
	"./databases"
	"./memqueue"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func loadEnv() {
	checkErr(godotenv.Load())
}

func printStart(url string) {
	fmt.Println("Starting request for: ", url)
}
func printComplete() {
	fmt.Printf("Complete.... Waiting %s seconds for next one.\n", os.Getenv("WAITING_TIME"))
	fmt.Println("-------------------------------------------")
	fmt.Println("")
}

func mapToSlice(m map[string]int) []string {
	var slice []string
	for k := range m {
		slice = append(slice, k)
	}
	return slice
}

func checkRequestError(err error, dbManager databases.DataManager, next string, queue *memqueue.Queue) bool {
	if err != nil {
		log.Printf("\nThis URL has something wrong, better skip\n")
		go func() {
			dbManager.RemoveFromMysql(GetCleanUrl(next))
		}()

		/*Add done URL to last URL Memory Queue*/
		queue.Push(memqueue.URL{
			Domain: GetDomainFromUrl(next),
			Clean:  GetCleanUrl(next),
		})
		return true
	}
	return false
}
