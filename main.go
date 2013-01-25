package main

import (
	"flag"
	"fmt"
	"gpio"
	"gpio/scheduler"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	ON  = gpio.HIGH
	OFF = gpio.LOW
)

var S *scheduler.Scheduler

var gVerbose *bool

func log(msg string, args ...interface{}) {
	if *gVerbose {
		fmt.Printf(msg, args...)
	}
}

func main() {
	schedule := flag.String("s", "db/schedule.json", "\tSchedule, list of events")
	gVerbose = flag.Bool("v", false, "\t\tDisplay verbose debug messages.")
	random := flag.Uint("r", 0, "\t\tLoad n random events for the schedule")
	flag.Parse()

	log("Creating a new event scheduler...")
	S = scheduler.New()
	defer S.CloseGPIOPins()
	log("done\n")

	log("Reading in the database...")
	if *random != 0 {
		rand.Seed(int64(time.Now().Second()))
		S.GenerateRandomEvents(int(*random))
		//err := C.SaveSchedule(*schedule)
		//if err != nil {
		//	return
		//}
	} else {
		err := S.LoadSchedule(*schedule)
		if err != nil {
			return
		}
	}
	log("done\n")

	log("Starting event queue manager\n")
	go S.ManageEventQueue()

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
