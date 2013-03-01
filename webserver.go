package main

import (
	"fmt"
	"html/template"
	//"io/ioutil"
	"net/http"
	"strings"
)

var lengths = map[string] int {
	"index":    len("/"),
	"event":    len("/event/"),
}

var templates = template.Must(template.ParseFiles(
	"template/header.html",
	"template/footer.html",
	"template/event.html",
	"template/event-view.html",
	"template/event-edit.html",
))

func renderTemplate(w http.ResponseWriter, title string) {
	err := templates.ExecuteTemplate(w, title + ".html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "header")
	fmt.Fprintf(w, "<h1>Welcome to the home page %d</h1>", lengths["index"])
	// list the buildings available
	// sublist the plans available
	renderTemplate(w, "footer")
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	args := strings.Split(r.URL.Path[lengths["event"]+1:], "/")
	title := ""
	if len(args) == 1 {
		title = args[0]
	}
	if len(args) > 1 {
		args = args[1:]
	}
	renderTemplate(w, "header")
	fmt.Fprintf(w, "Welcome to the event %s page", args)
	fmt.Fprintf(w, "Welcome to the event %s page", title)

	switch title {
	case "edit":
		renderTemplate(w, "event-edit")
	default:
		renderTemplate(w, "event")
	}
	renderTemplate(w, "footer")
}

/*
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
*/
