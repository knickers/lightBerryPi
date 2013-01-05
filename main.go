package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gpio"
	"io/ioutil"
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

var gVerbose *bool

func log(msg string, args ...interface{}) {
	if *gVerbose {
		fmt.Printf(msg, args...)
	}
}

func main() {
	schedule := flag.String("s", "schedule.json", "\tSchedule, list of events")
	gVerbose = flag.Bool("v", false, "\t\tDisplay verbose debug messages.")
	random := flag.Uint("r", 0, "\t\tLoad n random events for the schedule")
	flag.Parse()

	http.HandleFunc("/", webserv)
	fmt.Println("Serving on port: 8080")
	go http.ListenAndServe(":8080", nil)

	log("Creating a new light controller...")
	C := new(controller)
	defer C.closeGPIOPins()
	log("done\n")

	log("Reading in the database...")
	if *random != 0 {
		rand.Seed(int64(time.Now().Second()))
		C.generateRandomEvents(int(*random))
		//err := C.saveSchedule(*schedule)
		//if err != nil {
		//	return
		//}
	} else {
		err := C.loadSchedule(*schedule)
		if err != nil {
			return
		}
	}
	log("done\n")

	log("Starting event queue manager\n")
	go C.manageEventQueue()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	return
}
