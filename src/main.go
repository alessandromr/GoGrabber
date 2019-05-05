package main

import (
	"./Databases"
	"./MemQueue"
	"context"
	"os"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
)

func main() {

	/*
		Gui
	*/
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

	gridOpts, err := gridLayout(w) // equivalent to contLayout(w)
	if err != nil {
		panic(err)
	}

	if err := c.Update(rootID, gridOpts...); err != nil {
		panic(err)
	}

	/*
		GoGrabber
	*/
	loadEnv()
	//Prepare the Database clients
	var dbManager = databases.DataManager{}
	dbManager.SetClients()
	dbManager.RedisClient.FlushDB()

	//Prepare queues and channels
	queue := memqueue.Queue{Qlen: 1500}
	recent := memqueue.Queue{Qlen: 40}
	jobs := make(chan string, 100)

	jobs <- os.Getenv("START_URL")

	//Start routines
	go worker(jobs, dbManager, &queue, &recent, w)
	go inserter(jobs, dbManager, &queue, &recent)
	go supervisor(&recent, &queue, t)

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}
	if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quitter), termdash.RedrawInterval(redrawInterval)); err != nil {
		panic(err)
	}

}
