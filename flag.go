package main

import (
	"flag"
	"fmt"
)

type Config struct {
	AlertPercentage float64
	File            string
}

func GetConfig() (cfg Config, err error) {
	var (
		alertPercentage = flag.Float64("alert-percentage", 0.25, "Percentage above average to alert at")
		file            = flag.String("file", "", "File to tail")
	)
	flag.Parse()
	if *alertPercentage <= 0.0 {
		err = fmt.Errorf("alert-percentage must be greater than 0")
		return
	}
	if *file == "" {
		err = fmt.Errorf("file must be set")
		return
	}

	cfg = Config{
		AlertPercentage: *alertPercentage,
		File:            *file,
	}
	return
}
