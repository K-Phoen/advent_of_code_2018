package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type LogType byte

const NewShiftBegins LogType = 0
const GuardWakesUp LogType = 1
const GuardFallsAsleep LogType = 2

type Log struct {
	timestamp time.Time
	line      string
}

func (log *Log) logType() LogType {
	if strings.Index(log.line, "begins shift") != -1 {
		return NewShiftBegins
	} else if strings.Index(log.line, "wakes up") != -1 {
		return GuardWakesUp
	} else {
		return GuardFallsAsleep
	}
}

func (log *Log) extractGuardId() int {
	hash := strings.Index(log.line, "#")

	if hash == -1 {
		fmt.Fprintf(os.Stderr, "Expected to find guard ID in: '%s'\n", log.line)
		os.Exit(1)
	}

	space := strings.Index(log.line[hash:], " ")

	id, _ := strconv.Atoi(log.line[hash+1 : hash+space])

	return id
}

func extractTimestamp(line string) time.Time {
	end := strings.Index(line, "]")

	// yep, the errors are ignored. Because it's a game and I'm lazy.
	timestamp, _ := time.Parse("2006-01-02 15:04", line[1:end])

	return timestamp
}

func extractLogContent(line string) string {
	end := strings.Index(line, "]")

	return line[end+1:]
}

func readLogs(fileName string) []Log {
	var logs []Log

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open input file '%s'\n", fileName)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		logs = append(logs, Log{
			timestamp: extractTimestamp(line),
			line:      extractLogContent(line),
		})
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(2)
	}

	return logs
}
