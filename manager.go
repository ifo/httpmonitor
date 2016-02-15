package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

const (
	printInterval                = 10 * time.Second // how often to print results
	recentHistoryInterval        = 2 * time.Minute  // the length of recent history
	recentHistoryIntervalSeconds = uint64(recentHistoryInterval / time.Second)
)

func DataManager(in chan ProcessedLine) {
	// setup
	reqInfo := map[string]int{}
	var (
		totalHits     uint64 = 0
		twoMinuteHits uint64 = 0
	)

	// recur
	timer := make(chan struct{}, 1)
	go TenSecondTimer(timer)
	start := time.Now()
	for {
		select {
		case <-timer:
			fmt.Println(reqInfo)
			fmt.Println("avg hits/sec: ",
				totalHits/uint64(time.Since(start).Seconds()))
			fmt.Println("two-minute avg hits/sec: ",
				atomic.LoadUint64(&twoMinuteHits)/recentHistoryIntervalSeconds)
			fmt.Println("total hits: ", totalHits)
		case l := <-in:
			totalHits += 1
			atomic.AddUint64(&twoMinuteHits, 1)
			go DecrementAfterInterval(&twoMinuteHits)
			reqInfo[l.Section] += 1
		}
	}
}

func TenSecondTimer(timer chan<- struct{}) {
	for {
		time.Sleep(printInterval)
		timer <- struct{}{}
	}
}

func DecrementAfterInterval(v *uint64) {
	time.Sleep(recentHistoryInterval)
	atomic.AddUint64(v, ^uint64(0)) // subtracts 1
}
