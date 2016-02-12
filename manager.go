package main

import (
	"fmt"
	"sync/atomic"
	"time"
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

	for {
		select {
		case <-timer:
			fmt.Println(reqInfo)
			fmt.Println("total hits: ", totalHits)
			fmt.Println("two minute: ", atomic.LoadUint64(&twoMinuteHits))
		case l := <-in:
			totalHits += 1
			atomic.AddUint64(&twoMinuteHits, 1)
			go SubtractAfterTwoMinutes(&twoMinuteHits)
			reqInfo[l.Section] += 1
		}
	}
}

func TenSecondTimer(timer chan<- struct{}) {
	for {
		timer <- struct{}{}
		time.Sleep(10 * time.Second)
	}
}

func SubtractAfterTwoMinutes(v *uint64) {
	time.Sleep(2 * time.Minute)
	atomic.AddUint64(v, ^uint64(0)) // subtracts 1
}
