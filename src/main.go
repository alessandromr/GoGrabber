package main

import (
	"./databases"
	"./memqueue"
	"context"
	// "log"
	"os"
	// "time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
)

func main() {
	t, err := termbox.New(termbox.ColorMode(terminalapi.ColorMode256))
	if err != nil {
		panic(err)
	}
	defer t.Close()

	c, err := container.New(t, container.ID(rootID))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	w, err := newWidgets(ctx, c)
	if err != nil {
		panic(err)
	}
	// lb, err := newLayoutButtons(c, w)
	// if err != nil {
	// 	panic(err)
	// }
	// w.buttons = lb

	gridOpts, err := gridLayout(w) // equivalent to contLayout(w)
	if err != nil {
		panic(err)
	}

	if err := c.Update(rootID, gridOpts...); err != nil {
		panic(err)
	}

	loadEnv()
	//Prepare the Database clients
	var dbManager = databases.DataManager{}
	dbManager.SetClients()
	dbManager.RedisClient.FlushDB()

	queue := memqueue.Queue{Qlen: 1500}
	recent := memqueue.Queue{Qlen: 40}
	jobs := make(chan string, 100)

	jobs <- os.Getenv("START_URL")

	go worker(jobs, dbManager, &queue, &recent, w)
	go inserter(jobs, dbManager, &queue, &recent)

	// go func() {
	// 	for {
	// 		log.Println("")
	// 		log.Println("Status:")
	// 		log.Println("Queue len: ", len(queue.URLs))
	// 		log.Println("Pending jobs: ", len(jobs))
	// 		log.Println("")
	// 		time.Sleep(time.Second * 5)
	// 	}
	// }()

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}
	if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quitter), termdash.RedrawInterval(redrawInterval)); err != nil {
		panic(err)
	}

	// select {}

}
