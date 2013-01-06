package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gpio"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

type controller struct {
	pins   []gpio.Pin
	events []chan event // sorted by next closest event time
}

func (c *controller) exists(pin int) int {
	for i, p := range c.pins {
		if p.GetNumber() == uint(pin) {
			return i
		}
	}
	return -1
}

func (c *controller) setPinState(pin int, state gpio.State) {
	i := c.exists(pin)

	// This pin doesn't exist yet, create a new one
	if i == -1 {
		p, err := gpio.NewPin(uint(pin), gpio.OUTPUT)
		if err != nil {
			fmt.Println(err)
			return
		}
		c.pins = append(c.pins, *p)
		i = c.exists(pin)
	}

	c.pins[i].SetState(state)
}

func (c *controller) closeGPIOPins() {
	for _, p := range c.pins {
		p.Close()
	}
}

func (c *controller) pop() (event, error) {
	if len(c.events) < 1 {
		return event{}, errors.New("The events list is empty")
	}
	e := <-c.events[0]
	if len(c.events) > 1 {
		c.events = c.events[1:]
	} else {
		c.events = []chan event{}
	}
	return e, nil
}

func (c *controller) push(e event) error {
	evnt := make(chan event, 1)
	evnt <- e
	c.events = append(c.events, evnt)
	return nil
}

func (c *controller) insertInOrder(e event) error {
	ch := make(chan event, 1)
	ch <- e
	c.events = append(c.events, ch)
	for i := len(c.events) - 2; i >= 0; i-- {
		evnt := <-c.events[i]
		if e.NextTime.After(evnt.NextTime) {
			c.events[i] <- evnt
			break
		}
		c.events[i] <- evnt
		tmp := c.events[i]
		c.events[i] = c.events[i+1]
		c.events[i+1] = tmp
	}
	return nil
}

func (c *controller) getNextTime() (time.Time, error) {
	if len(c.events) == 0 {
		return time.Now(), errors.New("No events in the queue")
	}

	e := <-c.events[0]
	nextTime := e.NextTime
	c.events[0] <- e

	return nextTime, nil
}

func (c *controller) manageEventQueue() {
	log("Entering main loop\n")
	for {
		log("Popping off the next event\n")
		event, err := c.pop()
		if err == nil {
			now := time.Now()

			fmt.Printf("Sleeping %v till next event at %v...", event.NextTime.Sub(now), event.NextTime)
			time.Sleep(event.NextTime.Sub(now))
			fmt.Println("done")

			log("Setting the gpio pin states to %s.", event.State.String())
			for _, pin := range event.Pins {
				log("%d.", pin)
				/*
				err = C.setPinState(pin, event.state)
				if err != nil {
					break
				}
				*/
			}
			log("done\n")

			log("Updating the next time for this event...")
			err = event.updateNextTime()
			if err != nil {
				fmt.Println(err)
				break
			}
			log("done\n")

			log("Putting this event back on the queue...")
			err = c.insertInOrder(event)
			if err != nil {
				fmt.Println(err)
				break
			}
			log("done\n")
		} else {
			log(err.Error() + "\n")
			time.Sleep(time.Second)
		}
	}
}

func (c *controller) generateRandomEvents(num int) {
	for i := 0; i < num; i++ {
		// up to five pins per event
		n := rand.Int()%5 + 1
		var pins []int
		for j := 0; j < n; j++ {
			pins = append(pins, rand.Int()%8)
		}
		// on or off
		state := gpio.State(rand.Int() % 2)
		// up to twenty seconds in the future
		dur, err := time.ParseDuration(fmt.Sprintf("%ds", rand.Int()%20+1))
		if err != nil {
			fmt.Println(err)
			break
		}
		nextT := time.Now().Add(dur)
		// choose the days of the week to be applied
		var days []bool
		for j := 0; j < 7; j++ {
			r := false
			if rand.Int()%2 == 0 {
				r = true
			}
			days = append(days, r)
		}
		// choose the weeks of the year to be applied
		var weeks []bool
		for j := 0; j < 52; j++ {
			r := false
			if rand.Int()%2 == 0 {
				r = true
			}
			weeks = append(weeks, r)
		}
		c.insertInOrder(event{pins, state, nextT, days, weeks})
	}
}

func (c *controller) saveSchedule(file string) error {
	var events []event
	for _, e := range c.events {
		event := <-e
		events = append(events, event)
		e <- event
	}
	//bytes, err := json.Marshal(events)
	bytes, err := json.MarshalIndent(events, "", "\t")
	if err != nil {
		fmt.Println("Marshal:", err)
		return err
	}
	err = ioutil.WriteFile(file, bytes, os.FileMode(0664))
	if err != nil {
		fmt.Println("WriteOut:", err)
		return err
	}
	return nil
}

func (c *controller) loadSchedule(file string) error {
	fp, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("ReadFile:", err)
		return err
	}
	//log("%s\n", string(fp))
	var events []event
	err = json.Unmarshal(fp, &events)
	if err != nil {
		fmt.Println("Unmarshal:", err)
		return err
	}
	//log("%v\n", events)
	for _, e := range events {
		err = c.insertInOrder(e)
		if err != nil {
			fmt.Println("Insert:", err)
			return err
		}
	}
	return nil
}
