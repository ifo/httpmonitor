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

func (a Alert) Print() {
	if a.IsEmpty() {
		return
	}
	if a.IsCurrent() {
		fmt.Println("!!!=== WARNING: Alert is in progress ===!!!")
		a.Print()
		fmt.Println("!!!=========== current alert ===========!!!")
	} else {
		a.Print()
	}
}

func (a Alert) PrintRecovery() {
	if a.IsEmpty() {
		return
	}
	fmt.Println("!!!=== Alert recovered ===!!!")
	a.Print()
	fmt.Println("!!!=== alert has ended ===!!!")
}

// Empty alerts have not yet started
func (a Alert) IsEmpty() bool {
	return a.Start.IsZero()
}

// Current Alerts have not yet ended
func (a Alert) IsCurrent() bool {
	return !a.Start.IsZero() && a.End.IsZero()
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
		"Alert triggered at %s with hits/sec = (%.f), recovered at %s",
		a.Start.Format(time.UnixDate),
		a.HitsPerSec,
		a.End.Format(time.UnixDate),
	)
}
