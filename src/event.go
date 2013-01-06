package main

import (
	"errors"
	//"fmt"
	"gpio"
	"time"
)

type event struct {
	Pins        []int
	State       gpio.State
	NextTime    time.Time
	RepeatDays  []bool
	RepeatWeeks []bool
}

func (e *event) updateNextTime() error {
	// e.nextTime is thisTime right now
	today := e.NextTime.Weekday()
	_, thisWeek := e.NextTime.ISOWeek()

	firstTime := true
	for {
		day := 0
		if firstTime {
			day = int(today) + 1
			firstTime = false
		}
		for ; day < 7; day++ {
			// Add one more day to the wait time
			e.NextTime = e.NextTime.Add(24 * time.Hour)
			// If today is enabled then we're done
			if e.RepeatWeeks[(thisWeek-1)%52] && e.RepeatDays[day] {
				return nil
			}
		}
		thisWeek++
	}

	return errors.New("This point sould never be reached")
}

func (e *event) update(
	pins []int,
	state gpio.State,
	nextTime time.Time,
	repeatDays, repeatWeeks []bool) error {
	e.Pins = pins
	e.State = state
	e.NextTime = nextTime
	e.RepeatDays = repeatDays
	e.RepeatWeeks = repeatWeeks
	return nil
}
