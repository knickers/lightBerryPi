package main

import (
	"fmt"
	//"html/template"
	//"io/ioutil"
	"net/http"
)

var lengths = map[string] int {
	"index":    len("/"),
	"edit":     len("/edit/"),
	"plan":     len("/plan/"),
	"event":    len("/event/"),
	"login":    len("/login/"),
	"schedule": len("/schedule/"),
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page %d", lengths["index"])
	// list the buildings available
	// sublist the plans available
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[lengths["edit"]:]
	fmt.Fprintf(w, "Welcome to the edit %s page", title)
	// switch for what they are trying to edit
	switch title {
	case "building":
		fmt.Fprintf(w, "ho")
	case "floor":
		fmt.Fprintf(w, "hum")
	case "zone":
		fmt.Fprintf(w, "ha")
	case "event":
		fmt.Fprintf(w, "he")
	default:
		fmt.Fprintf(w, "no")
	}
	// building? floor? zone? event?
}

func floorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the floor plan page %d", lengths["plan"])
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the event page %d", lengths["event"])
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the login page %d", lengths["login"])
}

func scheduleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the schedule page %d", lengths["schedule"])
}
