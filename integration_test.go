package main

import (
	"testing"
	"time"
)

func Test_Alerting(t *testing.T) {
	timerChan := make(chan time.Time, 1)
	testChan := make(chan struct{}, 1)

	// make Config
	cfg := Config{
		AlertPercentage:       0.25,
		PrintTimer:            timerChan,
		RecentHistoryInterval: time.Nanosecond,
		GroupingResolution:    time.Nanosecond,
		Log:                   false,
		TestChannel:           testChan,
	}

	// setup State Manager, ignore input channel
	StateManager(cfg)

	// setup test states and stateReader
	testState := 0
	recentHits := []uint64{10, 100, 10}
	totalHits := []uint64{10, 110, 120}
	currentHitState := &HitState{}

	changeStateChan := make(chan bool, 1)
	sumTimeChan := make(chan bool, 1)
	doneChan := make(chan bool, 1)

	// modify testing functions
	ChangeState = func(hs *HitState) {
		currentHitState = hs                              // get state reference
		hs.StartTime = hs.StartTime.Add(-2 * time.Second) // increase duration by 2s
		hs.TotalHits = totalHits[testState]               // adjust total hits
		hs.RecentHits = 0                                 // zero for SumTimeGroup
		hs.RecentDurationSeconds = 2                      // set recent duration
		hs.AlertPercentage = cfg.AlertPercentage          // set alert percent
		changeStateChan <- true
	}
	SumTimeGroup = func(ts []TimeGroup) (out uint64) {
		out = recentHits[testState]
		sumTimeChan <- true
		return
	}
	Done = func() {
		doneChan <- true
	}

	testChan <- struct{}{}
	<-changeStateChan
	timerChan <- time.Now()
	<-sumTimeChan
	testState = 1
	<-doneChan

	if !currentHitState.TopAlert.IsEmpty() {
		t.Errorf("TopAlert should be empty, instead is %v",
			currentHitState.TopAlert)
	}

	testChan <- struct{}{}
	<-changeStateChan
	timerChan <- time.Now()
	<-sumTimeChan
	testState = 2
	<-doneChan

	if !currentHitState.TopAlert.IsCurrent() {
		t.Errorf("TopAlert should be current, instead is %v",
			currentHitState.TopAlert)
	}

	testChan <- struct{}{}
	<-changeStateChan
	timerChan <- time.Now()
	<-sumTimeChan
	<-doneChan

	if currentHitState.TopAlert.End.IsZero() {
		t.Errorf("TopAlert should be ended, instead is %v",
			currentHitState.TopAlert)
	}

	testChan <- struct{}{}
	<-changeStateChan
	timerChan <- time.Now()
	<-sumTimeChan
	<-doneChan

	if !currentHitState.TopAlert.IsEmpty() {
		t.Errorf("TopAlert should be empty, instead is %v",
			currentHitState.TopAlert)
	}
	if len(currentHitState.PastAlerts) != 1 {
		t.Errorf("PastAlerts should contain one alert, instead has %v",
			currentHitState.PastAlerts)
	}
}
