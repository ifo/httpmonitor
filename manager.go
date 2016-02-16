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
		twoMinHits []TimeGroup
	)

	ticker := time.NewTicker(printInterval)
	start := time.Now()
	go func() {
		for {
			select {
			case <-ticker.C:
				state["total"] = totalHits
				twoMinHits = RemoveOldTimeGroups(twoMinHits)
				state["recent"] = SumTimeGroup(twoMinHits)
				LogState(state, start)
			case l := <-input:
				totalHits += 1
				twoMinHits = GroupBySecond(twoMinHits, l.Time)
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

// A TimeGroup keeps track of the number of hits that have occured within a
// particular time range. In this case, 1 second. This was added to
// greatly reduce memory requirements during large traffic spikes.
type TimeGroup struct {
	Time  time.Time
	Count uint64
}

func GroupBySecond(ts []TimeGroup, t time.Time) []TimeGroup {
	l := len(ts)
	if l > 0 && t.Sub(ts[l-1].Time) <= time.Second {
		ts[l-1].Count += 1
		return ts
	} else {
		return append(ts, TimeGroup{Time: t, Count: 1})
	}
}

func RemoveOldTimeGroups(ts []TimeGroup) []TimeGroup {
	n := time.Now().Add(-recentHistoryInterval)
	i := 0
	for ; i < len(ts); i++ {
		if ts[i].Time.After(n) {
			break
		}
	}
	return ts[i:]
}

func SumTimeGroup(ts []TimeGroup) (sum uint64) {
	for _, t := range ts {
		sum += t.Count
	}
	return
}
