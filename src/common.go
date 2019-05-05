package main

import (
	"./Databases"
	"./MemQueue"
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

func printStart(url string, w *widgets) {
	addText("Starting request for: "+url, w.doneText)
}
func printComplete(w *widgets) {
	addText("Complete.... Waiting "+os.Getenv("WAITING_TIME")+" seconds for next one.", w.doneText)
}

func mapToSlice(m map[string]int) []string {
	var slice []string
	for k := range m {
		slice = append(slice, k)
	}
	return slice
}

func checkRequestError(err error, dbManager databases.DataManager, next string, queue *memqueue.Queue, w *widgets) bool {
	if err != nil {
		addText("\nThis URL has something wrong, better skip it\n", w.errorsText)
		go func() {
			dbManager.RemoveFromMysql(GetCleanURL(next))
		}()

		/*Add done URL to last URL Memory Queue*/
		queue.Push(memqueue.URL{
			Domain: GetDomainFromURL(next),
			Clean:  GetCleanURL(next),
		})
		return true
	}
	return false
}
