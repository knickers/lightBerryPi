package main

import (
	"flag"
	"fmt"
	"gpio"
	"math/rand"
	"msg"
	"net/http"
	"os"
	"os/signal"
	"scheduler"
	"time"
)

const (
	ON  = gpio.HIGH
	OFF = gpio.LOW
)

type extention struct {
	Pins  []int
	State gpio.State
}

type program struct {
	sched *scheduler.Scheduler
	pins  []gpio.Pin
	ios   map[int]extention
}

func (P *program) hasPin(pin int) int {
	for i, p := range P.pins {
		if p.GetNumber() == uint(pin) {
			return i
		}
	}
	return -1
}

func (P *program) setPinState(pin int, state gpio.State) error {
	i := P.hasPin(pin)

	// If this pin doesn't exist yet, create a new one
	if i == -1 {
		p, err := gpio.NewPin(uint(pin), gpio.OUTPUT)
		if err != nil {
			fmt.Println(err)
			return err
		}
		P.pins = append(P.pins, *p)
		i = len(P.pins) - 1
	}

	P.pins[i].SetState(state)
	return nil
}

func (P *program) closePins() {
	for _, p := range P.pins {
		p.Close()
	}
}

func (P *program) randomHelper() func(*scheduler.Event) {
	return func(event *scheduler.Event) {
		msg.Log("    Random pins for: ", event)
		n := rand.Int()%8 + 1 // eight pins on the raspberry pi
		var pins []int
		for j := 0; j < n; j++ {
			pins = append(pins, rand.Int()%8)
		}
		st := gpio.State(rand.Int() % 2)

		msg.Logln(" State:", st.String(), "Pins:", pins)
		P.ios[event.Id()] = extention{pins, st}
		//event.Action = P.makeEventAction(pins, st)
	}
}

func (P *program) loadHelper() func(*scheduler.Event, interface{}) {
	return func(event *scheduler.Event, data interface{}) {
		/*
		pins := extention(data)
		P.ios[event.Id()] = pins
		event.Action = P.makeEventAction(pins.Pins, pins.State)
		*/
	}
}

func (P *program) saveHelper() func(scheduler.Event) interface{} {
	return func(event scheduler.Event) interface{} {
		return P.ios[event.Id()]
	}
}

func (P *program) makeEventAction(pins []int, state gpio.State) func() error {
	return func() error {
		for _, pin := range pins {
			err := P.setPinState(pin, state)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

var P *program

func main() {
	msg.Verbose = flag.Bool("v", false, "\t\tDisplay verbose debug messages.")
	schedule := flag.String("s", "db/schedule.json", "Schedule list of events")
	random := flag.Int("r", 0, "\t\tLoad n random events for the schedule")
	flag.Parse()

	msg.Log("Creating a new event scheduler...")
	P = new(program)
	P.sched = scheduler.New()
	P.ios = make(map[int]extention)
	defer P.closePins()
	msg.Logln("done")

	msg.Logln("Reading in the database...")
	if *random > 0 {
		rand.Seed(int64(time.Now().Second()))
		P.sched.GenerateRandomEvents(*random, P.randomHelper())
		//err := P.sched.SaveSchedule(*schedule, P.saveHelper())
		//if err != nil {
		//	return
		//}
	} else {
		err := P.sched.LoadSchedule(*schedule, P.loadHelper())
		if err != nil {
			return
		}
	}
	msg.Logln("done")

	go P.sched.ManageEventQueue()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/event/", eventHandler)
	go http.ListenAndServe(":8080", nil)
	fmt.Println("Serving on port: 8080")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	return
}
