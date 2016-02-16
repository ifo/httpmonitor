package main

import (
	"log"

	"github.com/hpcloud/tail"
)

// TODO cleanup const usage
const (
	numProcessors = 200
)

func main() {
	cfg, err := GetConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// tail file
	t, err := tail.TailFile(cfg.File, tail.Config{
		//Location: &tail.SeekInfo{Offset: 0, Whence: 2}, // start at end of file
		ReOpen: true,
		Follow: true,
		Logger: tail.DiscardingLogger,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// create state input channel
	input := StateManager(printInterval)

	// kickoff line process workers
	for i := 0; i < numProcessors; i++ {
		go LineProcessWorker(t.Lines, input)
	}

	// wait forever
	done := make(chan struct{}, 1)
	<-done
}
