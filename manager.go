package main

import (
	"fmt"
	"time"
)

func DataManager(in chan ProcessedLine) {
	// setup
	reqInfo := map[string]int{}

	// recur
	timer := make(chan struct{}, 1)
	go TenSecondTimer(timer)

	for {
		select {
		case <-timer:
			fmt.Println(reqInfo)
		case l := <-in:
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
