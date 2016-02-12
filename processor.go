package main

import (
	"strconv"
	"strings"
	"time"
)

type ProcessedLine struct {
	Method       string
	Section      string
	Protocol     string
	ResponseCode int
	IP           string
	Time         time.Time
}

func ProcessLogLine(l string, out chan<- ProcessedLine) error {
	if line, err := LineProcessor(l); err != nil {
		return err
	} else {
		out <- line
		return nil
	}
}

// LineProcessor assumes a common log format:
// https://www.w3.org/Daemon/User/Config/Logging.html#common-logfile-format
func LineProcessor(l string) (out ProcessedLine, err error) {
	//dateStart := strings.Index(l, "[")
	dateEnd := strings.Index(l, "]")

	//date := l[dateStart+1 : dateEnd]
	ipEnd := strings.Index(l, " ")
	ip := l[:ipEnd]

	queryParts := strings.Fields(l[dateEnd+2:])

	respCode, err := strconv.Atoi(queryParts[3])
	if err != nil {
		return
	}

	out = ProcessedLine{
		Method:       queryParts[0][1:], // remove " at the start
		Section:      queryParts[1],
		Protocol:     queryParts[2][:len(queryParts[2])-1], // remove " at the end
		ResponseCode: respCode,
		IP:           ip,
		Time:         time.Now().UTC(),
	}
	return
}
