package main

import (
	"flag"
	"fmt"
	"gpio"
	"gpio/scheduler"
	"math/rand"
	"msg"
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

func main() {
	msg.Verbose = flag.Bool("v", false, "\t\tDisplay verbose debug messages.")
	schedule := flag.String("s", "db/schedule.json", "Schedule list of events")
	random := flag.Int("r", 0, "\t\tLoad n random events for the schedule")
	flag.Parse()

	msg.Log("Creating a new event scheduler...")
	S = scheduler.New()
	defer S.CloseGPIOPins()
	msg.Log("done\n")

	msg.Log("Reading in the database...")
	if *random > 0 {
		rand.Seed(int64(time.Now().Second()))
		S.GenerateRandomEvents(*random)
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
	msg.Log("done\n")

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
