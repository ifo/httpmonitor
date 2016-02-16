package main

import (
	"fmt"
	"time"
)

// TODO cleanup const usage
const (
	printInterval                = 10 * time.Second // how often to print results
	recentHistoryInterval        = 2 * time.Minute  // the length of recent history
	recentHistoryIntervalSeconds = uint64(recentHistoryInterval / time.Second)
)

func StateManager(printInterval time.Duration) chan<- ProcessedLine {
	input := make(chan ProcessedLine)
	state := map[string]uint64{}
	var (
		totalHits  uint64 = 0
		twoMinHits []time.Time
	)

	ticker := time.NewTicker(printInterval)
	start := time.Now()
	go func() {
		for {
			select {
			case <-ticker.C:
				state["total"] = totalHits
				twoMinHits = RemoveOld(twoMinHits)
				state["recent"] = uint64(len(twoMinHits))
				LogState(state, start)
			case l := <-input:
				totalHits += 1
				twoMinHits = append(twoMinHits, l.Time)
				state[l.Section] += 1
			}
		}
	}()
	return input
}

func LogState(s map[string]uint64, start time.Time) {
	fmt.Println(s)
	fmt.Println("avg hits/sec: ", s["total"]/uint64(time.Since(start).Seconds()))
	fmt.Println("two-minute avg hits/sec: ", s["recent"]/recentHistoryIntervalSeconds)
	fmt.Println("total hits: ", s["total"])
}

func RemoveOld(ts []time.Time) []time.Time {
	n := time.Now().Add(-recentHistoryInterval)
	i := 0
	for ; i < len(ts); i++ {
		if ts[i].After(n) {
			break
		}
	}
	return ts[i:]
}
