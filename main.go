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

func (P *program) randomHelper(e scheduler.Event) {
	// up to eight pins per event
	n := rand.Int()%8 + 1
	var pins []int
	for j := 0; j < n; j++ {
		pins = append(pins, rand.Int()%8)
	}
	st := gpio.State(rand.Int() % 2)
	P.ios[e.Id()] = extention{pins, st}
}

func (P *program) loadHelper(e scheduler.Event) {
}

func (P *program) saveHelper(e scheduler.Event) {
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
	defer P.closePins()
	msg.Log("done\n")

	msg.Log("Reading in the database...")
	if *random > 0 {
		rand.Seed(int64(time.Now().Second()))
		P.sched.GenerateRandomEvents(*random, nil/*P.randomHelper*/)
		//err := P.sched.SaveSchedule(*schedule, P.saveHelper)
		//if err != nil {
		//	return
		//}
	} else {
		err := P.sched.LoadSchedule(*schedule, nil/*P.loadHelper*/)
		if err != nil {
			return
		}
	}
	msg.Log("done\n")

	go P.sched.ManageEventQueue()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/event/", eventHandler)
	http.HandleFunc("/floor/", floorHandler)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/schedule/", scheduleHandler)
	go http.ListenAndServe(":8080", nil)
	fmt.Println("Serving on port: 8080")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	return
}
