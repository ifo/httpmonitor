package main

import (
	"flag"
	"fmt"
	"time"
)

type Config struct {
	AlertPercentage       float64
	File                  string
	WorkerCount           int
	PrintInterval         time.Duration
	RecentHistoryInterval time.Duration
	GroupingResolution    time.Duration
}

func GetConfig() (cfg Config, err error) {
	var (
		alertPercentage = flag.Float64("alert-percentage", 0.25,
			"Percentage above average to alert at")
		file = flag.String("file", "",
			"File to tail")
		workerCount = flag.Int("worker-count", 200,
			"Number of workers (~200 seems to work best)")
		printInterval = flag.Duration("print", 10*time.Second,
			"Time between printingi (e.g. 10s)")
		recentHistoryInterval = flag.Duration("recent-history", 2*time.Minute,
			"Length of recent history (e.g. 2m)")
		groupingResolution = flag.Duration("grouping", 1*time.Second,
			"Resolution of the recent history interval (smaller uses more memory)")
	)
	flag.Parse()
	if *alertPercentage <= 0.0 {
		err = fmt.Errorf("alert-percentage must be greater than 0")
		return
	}
	if *file == "" {
		err = fmt.Errorf("file must be set")
		return
	}
	if *workerCount <= 0 {
		err = fmt.Errorf("worker-count must be greater than 0")
		return
	}
	if *printInterval < 1*time.Second {
		err = fmt.Errorf("print must be at least 1 second (1s)")
		return
	}
	if *recentHistoryInterval <= *printInterval {
		err = fmt.Errorf(
			"recent-history must be greater than print (print default is 10s)")
		return
	}
	if *groupingResolution < 1*time.Millisecond {
		err = fmt.Errorf("grouping must be at least 1 millisecond")
		return
	}

	cfg = Config{
		AlertPercentage:       *alertPercentage,
		File:                  *file,
		WorkerCount:           *workerCount,
		PrintInterval:         *printInterval,
		RecentHistoryInterval: *recentHistoryInterval,
		GroupingResolution:    *groupingResolution,
	}
	return
}
