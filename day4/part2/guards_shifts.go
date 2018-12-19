package main

import (
	"fmt"
	"math/bits"
	"os"
	"sort"
	"time"
)

type DailyActivity struct {
	timestamp time.Time
	detail    uint64
}

type GuardActivity struct {
	id      int
	records []DailyActivity
}

type ActivityBook struct {
	guards map[int]*GuardActivity
}

func (activity *DailyActivity) timeAsleep() int {
	return bits.OnesCount64(activity.detail)
}

func (activity *DailyActivity) registerSleepingPeriod(from time.Time, to time.Time) {
	var mask uint64 = 1
	for i := 0; i < 60; i++ {
		if i >= from.Minute() && i < to.Minute() {
			activity.detail = activity.detail | mask
		}

		mask = mask << 1
	}
}

func (activity *DailyActivity) Output() {
	fmt.Printf("%02d-%02d\t", activity.timestamp.Month(), activity.timestamp.Day())

	var mask uint64 = 1
	for i := 0; i < 60; i++ {
		b := activity.detail & mask

		if b != 0 {
			fmt.Print("#")
		} else {
			fmt.Print(".")
		}

		mask = mask << 1
	}

	fmt.Printf(" (asleep for %d minutes)\n", activity.timeAsleep())
}

func (activity *GuardActivity) Output() {
	for _, record := range activity.records {
		record.Output()
	}
}

func (activity *GuardActivity) timeAsleep() int {
	total := 0

	for _, record := range activity.records {
		total += record.timeAsleep()
	}

	return total
}

func (activity *DailyActivity) asleepMinutes() []int {
	var minutes []int

	var mask uint64 = 1
	for i := 0; i < 60; i++ {
		b := activity.detail & mask

		if b != 0 {
			minutes = append(minutes, i)
		}

		mask = mask << 1
	}

	return minutes
}

func (book *ActivityBook) registerGuard(id int) {
	if book.hasGuard(id) {
		return
	}

	book.guards[id] = &GuardActivity{
		id: id,
	}
}

func (book *ActivityBook) hasGuard(id int) bool {
	_, exists := book.guards[id]

	return exists
}

func (book *ActivityBook) addDailyActivityFor(id int, activity DailyActivity) {
	book.guards[id].records = append(book.guards[id].records, activity)
}

func (book *ActivityBook) Output() {
	for guardId, activity := range book.guards {
		fmt.Printf("Activity for gard #%d (asleep for %d minutes in total)\n", guardId, activity.timeAsleep())
		activity.Output()
	}
}

func (book *ActivityBook) FindMostAsleepGuard() int {
	guardId := 0
	max := 0

	for id, activity := range book.guards {
		if activity.timeAsleep() > max {
			guardId = id
			max = activity.timeAsleep()
		}
	}

	return guardId
}

func (book *ActivityBook) FindMostAsleepMinute(guardId int) (int, int) {
	activity := book.guards[guardId]
	asleepMap := make(map[int]int)
	total := 0
	minute := -1

	for _, record := range activity.records {
		for _, m := range record.asleepMinutes() {
			asleepMap[m] += 1

			if asleepMap[m] > total {
				total = asleepMap[m]
				minute = m
			}
		}
	}

	return minute, total
}

func (book *ActivityBook) FindMostFrequentlyAsleepOnSameMinute() (int, int, int) {
	id := -1
	count := 0
	minute := -1

	for guardId := range book.guards {
		m, c := book.FindMostAsleepMinute(guardId)

		if c > count {
			count = c
			id = guardId
			minute = m
		}
	}

	return id, minute, count
}

func ActivityBookFromLogs(logs []Log) *ActivityBook {
	book := &ActivityBook{
		guards: make(map[int]*GuardActivity),
	}

	sort.Slice(logs, func(i, j int) bool {
		return logs[i].timestamp.Before(logs[j].timestamp)
	})

	guardId := 0
	var dailyActivity DailyActivity
	var asleepSince time.Time

	for _, log := range logs {
		if log.logType() == NewShiftBegins && guardId != 0 {
			book.addDailyActivityFor(guardId, dailyActivity)
		}

		switch log.logType() {
		case NewShiftBegins:
			guardId = log.extractGuardId()

			book.registerGuard(guardId)
			dailyActivity = DailyActivity{
				timestamp: log.timestamp,
			}
		case GuardFallsAsleep:
			asleepSince = log.timestamp
		case GuardWakesUp:
			dailyActivity.registerSleepingPeriod(asleepSince, log.timestamp)
		}

		fmt.Printf("%v → %s (%d)\n", log.timestamp, log.line, guardId)
	}

	book.addDailyActivityFor(guardId, dailyActivity)

	return book
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "FILE")
		os.Exit(1)
	}

	logs := readLogs(os.Args[1])

	activityBook := ActivityBookFromLogs(logs)
	fmt.Println("-------------------")

	activityBook.Output()

	guardId, minute, count := activityBook.FindMostFrequentlyAsleepOnSameMinute()
	fmt.Printf("Guard #%d spent minute %d asleep more than any other guard or minute (%d times)\n", guardId, minute, count)

	fmt.Printf("Answer: %d * %d = %d\n", guardId, minute, guardId*minute)
}
