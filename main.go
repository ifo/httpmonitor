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

	t, err := tail.TailFile(cfg.File, tail.Config{
		//Location: &tail.SeekInfo{Offset: 0, Whence: 2}, // start at end of file
		ReOpen: true,
		Follow: true,
		Logger: tail.DiscardingLogger,
	})
	if err != nil {
		log.Fatalln(err)
	}

	ch := make(chan ProcessedLine)

	go DataManager(ch)

	for line := range t.Lines {
		go ProcessLogLine(line.Text, ch)
	}
}
