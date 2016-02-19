package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/hpcloud/tail"
)

// ProcessedLine is produced by the LineProcessWorker
// and consumed by the StateManager
type ProcessedLine struct {
	Method       string
	Path         string
	Section      string
	Protocol     string
	ResponseCode int
	IP           string
	Time         time.Time
}

// LineProcessor assumes a common log format, and that
// rfc931 and authuser don't contain the string `] "`
//
// Common Log Format:
// https://www.w3.org/Daemon/User/Config/Logging.html#common-logfile-format
//
// Common Log Format looks like:
// remotehost rfc931 authuser [date] "request" status bytes
// where "request" is: "METHOD /path PROTOCOL"
func LineProcessor(line string) (out ProcessedLine, err error) {
	ipEnd := strings.Index(line, " ")
	ip := line[:ipEnd]

	requestBegin := strings.Index(line, `] "`) + 3 // skip to start with METHOD
	requestParts := strings.Fields(line[requestBegin:])

	respCode, err := strconv.Atoi(requestParts[3])
	if err != nil {
		return
	}

	section := requestParts[1]
	sectionEnd := strings.Index(section[1:], "/")
	if sectionEnd != -1 {
		section = section[:sectionEnd+1]
	}

	out = ProcessedLine{
		Method:       requestParts[0],
		Path:         requestParts[1],
		Section:      section,
		Protocol:     requestParts[2][:len(requestParts[2])-1], // remove trailing "
		ResponseCode: respCode,
		IP:           ip,
		Time:         time.Now(),
	}
	return
}

// LineProcessWorker takes a tail line and creates a ProcessedLine struct
// It then sends that to the StateManager's input channel
func LineProcessWorker(in <-chan *tail.Line, stats chan<- ProcessedLine) {
	for l := range in {
		line, err := LineProcessor(l.Text)
		if err == nil {
			stats <- line
		}
	}
}
