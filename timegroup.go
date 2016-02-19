package main

import (
	"time"
)

// A TimeGroup keeps track of the number of hits that have occured within a
// particular time range. In this case, Config.GroupingResolution. This was
// added to greatly reduce memory requirements during large traffic spikes.
type TimeGroup struct {
	Time  time.Time
	Count float64
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

// SumTimeGroup is assignable for testing
var SumTimeGroup = sumTimeGroup

func sumTimeGroup(ts []TimeGroup) (sum float64) {
	for _, t := range ts {
		sum += t.Count
	}
	return
}
