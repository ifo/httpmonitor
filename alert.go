package main

import (
	"fmt"
	"time"
)

type Alert struct {
	Start      time.Time
	HitsPerSec float64
	End        time.Time
}

// Empty alerts have not yet started
func (a Alert) IsEmpty() bool {
	return a.Start.IsZero()
}

// Current Alerts have not yet ended
func (a Alert) IsCurrent() bool {
	return !a.Start.IsZero() && a.End.IsZero()
}

func (a Alert) MaxHitsPerSec(hps float64) float64 {
	if a.HitsPerSec < hps {
		return hps
	}
	return a.HitsPerSec
}

func (a Alert) Print() {
	if a.IsEmpty() {
		return
	}
	if a.IsCurrent() {
		fmt.Println("!!!=== WARNING: Alert is in progress ===!!!")
		fmt.Println("!!!=== ", a.AlertMessage())
		fmt.Println("!!!=========== current alert ===========!!!")
	} else {
		fmt.Println("|=== ", a.AlertMessage())
	}
}

func (a Alert) PrintRecovery() {
	if a.IsEmpty() {
		return
	}
	fmt.Println("!!!=== Alert recovered ===!!!")
	fmt.Println("!!!=== ", a.AlertMessage())
	fmt.Println("!!!=== alert has ended ===!!!")
}

func (a Alert) AlertMessage() string {
	if a.IsEmpty() {
		return ""
	}
	if a.IsCurrent() {
		return fmt.Sprintf(
			"High traffic generated an alert - hits/sec = %.f, triggered at %s",
			a.HitsPerSec,
			a.Start.Format(time.UnixDate),
		)
	}
	return fmt.Sprintf(
		"Recovered: Alert at %s with hits/sec = (%.f), recovered at %s",
		a.Start.Format(time.UnixDate),
		a.HitsPerSec,
		a.End.Format(time.UnixDate),
	)
}
