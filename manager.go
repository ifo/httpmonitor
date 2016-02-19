package main

import (
	"fmt"
	"time"
)

func StateManager(cfg Config) chan<- ProcessedLine {
	input := make(chan ProcessedLine)
	hs := HitState{
		HitMap:                map[string]uint64{},
		TotalHits:             0,
		RecentHits:            0,
		PastAlerts:            []Alert{},
		TopAlert:              Alert{},
		StartTime:             time.Now(),
		RecentDurationSeconds: uint64(cfg.RecentHistoryInterval / time.Second),
		AlertPercentage:       cfg.AlertPercentage,
	}
	hitsGroup := []TimeGroup{}

	ticker := time.NewTicker(cfg.PrintInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				hitsGroup = RemoveOldTimeGroups(hitsGroup, cfg.RecentHistoryInterval)
				hs = hs.Update(SumTimeGroup(hitsGroup))
				if cfg.Log {
					hs.Print()
				}
			case l := <-input:
				hs.TotalHits += 1
				hs.HitMap[l.Section] += 1
				hitsGroup = GroupByResolution(hitsGroup, l.Time,
					cfg.GroupingResolution)
			}
		}
	}()
	return input
}

// HitState contains all af the relevant logging data
type HitState struct {
	HitMap                map[string]uint64
	TotalHits             uint64
	RecentHits            uint64
	PastAlerts            []Alert
	TopAlert              Alert
	StartTime             time.Time
	RecentDurationSeconds uint64
	AlertPercentage       float64
}

// HitState.Update is a state machine with 3 possible states
// Empty Alert     -> Check for Alert (if so, move to Active Alert)
// Active Alert    -> Check for end of Alert (if so, move to Recovered Alert)
// Recovered Alert -> Save Alert in PastAlerts, move to Empty Alert
func (hs HitState) Update(recentHits uint64) HitState {
	hs.RecentHits = recentHits
	switch {
	// Empty Alert
	case hs.TopAlert.IsEmpty():
		hs.TopAlert = hs.CheckForAlert()
	// Active Alert
	case hs.TopAlert.IsCurrent():
		hs.TopAlert.HitsPerSec = hs.TopAlert.MaxHitsPerSec(
			float64(hs.RecentHits) / float64(hs.RecentDurationSeconds))
		hs.TopAlert = hs.CheckForAlertRecovery()
	// Recovered Alert
	case !hs.TopAlert.IsEmpty() && !hs.TopAlert.IsCurrent():
		hs.PastAlerts = append(hs.PastAlerts, hs.TopAlert)
		hs.TopAlert = Alert{}
	}
	return hs
}

func (hs HitState) Print() {
	// Print Current Alert
	if !hs.TopAlert.IsEmpty() {
		if hs.TopAlert.IsCurrent() {
			hs.TopAlert.Print()
		} else {
			hs.TopAlert.PrintRecovery()
		}
		fmt.Println("")
	}

	// Print Past Alerts
	if len(hs.PastAlerts) > 0 {
		fmt.Println("|==== Past alerts ====|")
		for _, a := range hs.PastAlerts {
			a.Print()
		}
		fmt.Println("|=== end of alerts ===|")
		fmt.Println("")
	}

	// Print Stats
	// TODO improve print stats
	fmt.Println(hs.HitMap)
	fmt.Println("")

	// Print Interesting Facts
	// TODO improve interesting facts
	fmt.Println("=== Did you know? ===")
	fmt.Println("=== total hits: ", hs.TotalHits)
	fmt.Println("=== avg hits/sec: ", hs.TotalHits/uint64(time.Since(hs.StartTime).Seconds()))
	fmt.Println("=== recent hits/sec: ", hs.RecentHits/hs.RecentDurationSeconds)
	fmt.Println("")
}

func (hs HitState) CheckForAlert() (alert Alert) {
	tot64 := float64(hs.TotalHits) / float64(time.Since(hs.StartTime).Seconds())
	rec64 := float64(hs.RecentHits) / float64(hs.RecentDurationSeconds)
	boundary := tot64 + tot64*hs.AlertPercentage
	if rec64 > boundary {
		alert.Start = time.Now()
		alert.HitsPerSec = alert.MaxHitsPerSec(rec64)
	}
	return
}

func (hs HitState) CheckForAlertRecovery() (alert Alert) {
	alert = hs.TopAlert
	tot64 := float64(hs.TotalHits) / float64(time.Since(hs.StartTime).Seconds())
	rec64 := float64(hs.RecentHits) / float64(hs.RecentDurationSeconds)
	boundary := tot64 + tot64*hs.AlertPercentage
	if rec64 < boundary {
		alert.End = time.Now()
		alert.HitsPerSec = alert.MaxHitsPerSec(rec64)
	}
	return
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
