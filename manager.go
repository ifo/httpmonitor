package main

import (
	"fmt"
	"time"
)

func StateManager(cfg Config) chan<- ProcessedLine {
	input := make(chan ProcessedLine)
	hs := HitState{
		HitMap:                map[string]int{},
		TotalHits:             0,
		RecentHits:            0,
		PastAlerts:            []Alert{},
		TopAlert:              Alert{},
		StartTime:             time.Now(),
		RecentDurationSeconds: float64(cfg.RecentHistoryInterval / time.Second),
		AlertPercentage:       cfg.AlertPercentage,
	}
	hitsGroup := []TimeGroup{}

	go func() {
		for {
			select {
			case <-cfg.PrintTimer:
				hitsGroup = RemoveOldTimeGroups(hitsGroup, cfg.RecentHistoryInterval)
				hs = hs.Update(SumTimeGroup(hitsGroup))
				if cfg.Log {
					hs.Print()
				}
				Done() // Only used in testing
			case l := <-input:
				hs.TotalHits += 1
				hs.HitMap[l.Section] += 1
				hitsGroup = GroupByResolution(hitsGroup, l.Time,
					cfg.GroupingResolution)
			case <-cfg.TestChannel:
				ChangeState(&hs)
			}
		}
	}()
	return input
}

// ChangeState and Done are noops used for testing
var ChangeState = func(hs *HitState) {}
var Done = func() {}

// HitState contains all af the relevant logging data
type HitState struct {
	HitMap                map[string]int
	TotalHits             float64
	RecentHits            float64
	PastAlerts            []Alert
	TopAlert              Alert
	StartTime             time.Time
	RecentDurationSeconds float64
	AlertPercentage       float64
}

// HitState.Update is a state machine with 3 possible states
// Empty Alert     -> Check for Alert (if so, move to Active Alert)
// Active Alert    -> Check for end of Alert (if so, move to Recovered Alert)
// Recovered Alert -> Save Alert in PastAlerts, move to Empty Alert
func (hs HitState) Update(recentHits float64) HitState {
	hs.RecentHits = recentHits
	switch {
	// Empty Alert
	case hs.TopAlert.IsEmpty():
		hs.TopAlert = hs.CheckForAlert()
	// Active Alert
	case hs.TopAlert.IsCurrent():
		hs.TopAlert.HitsPerSec = hs.TopAlert.MaxHitsPerSec(
			hs.RecentHits / hs.RecentDurationSeconds)
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
	fmt.Printf("=== total hits: %.f\n", hs.TotalHits)
	fmt.Printf("=== avg hits/sec: %.f\n", hs.TotalHits/timeSince(hs.StartTime))
	fmt.Printf("=== recent hits/sec: %.f\n", hs.RecentHits/hs.RecentDurationSeconds)
	fmt.Println("")
}

func (hs HitState) CheckForAlert() (alert Alert) {
	tot64 := hs.TotalHits / timeSince(hs.StartTime)
	rec64 := hs.RecentHits / hs.RecentDurationSeconds
	boundary := tot64 + tot64*hs.AlertPercentage
	if rec64 > boundary {
		alert.Start = time.Now()
		alert.HitsPerSec = alert.MaxHitsPerSec(rec64)
	}
	return
}

func (hs HitState) CheckForAlertRecovery() (alert Alert) {
	alert = hs.TopAlert
	tot64 := hs.TotalHits / timeSince(hs.StartTime)
	rec64 := hs.RecentHits / hs.RecentDurationSeconds
	boundary := tot64 + tot64*hs.AlertPercentage
	if rec64 < boundary {
		alert.End = time.Now()
		alert.HitsPerSec = alert.MaxHitsPerSec(rec64)
	}
	return
}

func timeSince(t time.Time) float64 {
	return float64(time.Since(t).Seconds())
}
