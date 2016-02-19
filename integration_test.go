package main

import (
	"testing"
	"time"
)

func Test_Alerting(t *testing.T) {
	// setup timer and config
	timerChan := make(chan time.Time, 1)
	cfg := Config{
		AlertPercentage:       0.25,
		PrintTimer:            timerChan,
		RecentHistoryInterval: time.Nanosecond,
		GroupingResolution:    time.Nanosecond,
		Log:                   false,
		TestChannel:           make(chan bool, 1),
	}

	// setup State Manager, ignore input channel
	StateManager(cfg)

	// setup test state and create state step function
	stepToState := setup(cfg, timerChan)

	emptyHitState := stepToState(1)
	if !emptyHitState.TopAlert.IsEmpty() {
		t.Errorf("TopAlert should be empty, instead is %v", emptyHitState.TopAlert)
	}

	alertHitState := stepToState(2)
	if !alertHitState.TopAlert.IsCurrent() {
		t.Errorf("TopAlert should be current, instead is %v", alertHitState.TopAlert)
	}

	recoverHitState := stepToState(3)
	if recoverHitState.TopAlert.End.IsZero() {
		t.Errorf("TopAlert should be ended, instead is %v", recoverHitState.TopAlert)
	}

	finalHitState := stepToState(4)
	if !finalHitState.TopAlert.IsEmpty() {
		t.Errorf("TopAlert should be empty, instead is %v", finalHitState.TopAlert)
	}
	if len(finalHitState.PastAlerts) != 1 {
		t.Errorf("PastAlerts should contain one alert, instead has %v",
			finalHitState.PastAlerts)
	}
}

// setup returns a function that will step to the given state as set by
// recentHits and totalHits, and give the HitState output at that time
func setup(cfg Config, timerChan chan time.Time) func(int) HitState {
	extHitState := &HitState{}
	testState := 0
	recentHits := []float64{10, 100, 10, 10, 10}
	totalHits := []float64{10, 110, 120, 130, 140}

	changeState := make(chan bool, 1)
	sumTimeGroup := make(chan bool, 1)
	done := make(chan bool, 1)

	ChangeState = func(hs *HitState) {
		extHitState = hs // get state reference

		hs.StartTime = hs.StartTime.Add(-2 * time.Second) // increase duration by 2s
		hs.TotalHits = totalHits[testState]
		hs.RecentHits = 0 // set later by SumTimeGroup
		hs.RecentDurationSeconds = 2
		hs.AlertPercentage = cfg.AlertPercentage

		changeState <- true
	}
	SumTimeGroup = func(ts []TimeGroup) (out float64) {
		out = recentHits[testState]
		sumTimeGroup <- true
		return
	}
	Done = func() {
		done <- true
	}

	return func(nextState int) HitState {
		cfg.TestChannel <- true // kickoff ChangeState
		<-changeState           // hold until ChangeState is complete

		timerChan <- time.Now() // kickoff Print Timer HitGroup update
		<-sumTimeGroup          // hold until HitGroup has been updated

		testState = nextState // change to next state
		<-done                // hold until testState has been updated

		return *extHitState
	}
}
