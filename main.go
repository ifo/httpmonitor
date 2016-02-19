package main

import (
	"log"

	"github.com/hpcloud/tail"
)

func main() {
	cfg, err := GetConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// create tail file channel
	tailChan, err := tail.TailFile(cfg.File, tail.Config{
		Location: &tail.SeekInfo{Offset: 0, Whence: 2}, // start at end of file
		ReOpen:   true,
		Follow:   true,
		Logger:   tail.DiscardingLogger, // don't print tailing information
	})
	if err != nil {
		log.Fatalln(err)
	}

	// create state input channel
	input := StateManager(cfg)

	// kickoff line process workers
	for i := 0; i < cfg.WorkerCount; i++ {
		go LineProcessWorker(tailChan.Lines, input)
	}

	// wait forever
	done := make(chan struct{}, 1)
	<-done
}
