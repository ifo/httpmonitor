package main

import (
	"fmt"
	"time"
)

func StateManager(cfg Config) chan<- ProcessedLine {
	input := make(chan ProcessedLine)
	state := map[string]uint64{}
	var (
		totalHits  uint64 = 0
		recentHits []TimeGroup
	)

	ticker := time.NewTicker(cfg.PrintInterval)
	start := time.Now()
	go func() {
		for {
			select {
			case <-ticker.C:
				state["total"] = totalHits
				recentHits = RemoveOldTimeGroups(recentHits, cfg.RecentHistoryInterval)
				state["recent"] = SumTimeGroup(recentHits)
				LogState(state, start, uint64(cfg.RecentHistoryInterval/time.Second))
			case l := <-input:
				totalHits += 1
				recentHits = GroupByResolution(recentHits, l.Time,
					cfg.GroupingResolution)
				state[l.Section] += 1
			}
		}
	}()
	return input
}

func LogState(s map[string]uint64, start time.Time, intervalSeconds uint64) {
	fmt.Println(s)
	fmt.Println("avg hits/sec: ", s["total"]/uint64(time.Since(start).Seconds()))
	fmt.Println("recent-history avg hits/sec: ",
		s["recent"]/intervalSeconds)
	fmt.Println("total hits: ", s["total"])
}

// A TimeGroup keeps track of the number of hits that have occured within a
// particular time range. In this case, Config.GroupingResolution. This was
// added to greatly reduce memory requirements during large traffic spikes.
type TimeGroup struct {
	Time  time.Time
	Count uint64
}

func GroupByResolution(ts []TimeGroup, t time.Time,
	resolution time.Duration) []TimeGroup {
	l := len(ts)
	if l > 0 && t.Sub(ts[l-1].Time) <= resolution {
		ts[l-1].Count += 1
		return ts
	} else {
		return append(ts, TimeGroup{Time: t, Count: 1})
	}
}

func RemoveOldTimeGroups(ts []TimeGroup, interval time.Duration) []TimeGroup {
	n := time.Now().Add(-interval)
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
