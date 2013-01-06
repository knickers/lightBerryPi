package main

import (
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
)

var HEADER, _ = template.ParseFiles("header.html")
var FOOTER, _ = template.ParseFiles("footer.html")

func (c *controller) webserv(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		f, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "create: " + err.Error(), 500)
			return
		}
		defer f.Close()
		t, err := ioutil.TempFile(".", "image-")
		if err != nil {
			http.Error(w, "temp: " + err.Error(), 500)
			return
		}
		defer t.Close()
		if _, err := io.Copy(t, f); err != nil {
			http.Error(w, "copy: " + err.Error(), 500)
			return
		}
		http.Redirect(w, r, "/view?id=" + t.Name()[6:], 302)
	} else {
		page, err := template.ParseFiles("template/event-list.go.html")
		if err != nil {
			http.Error(nil, "Page Parse Error: " + err.Error(), 500)
			return
		}
		HEADER.Execute(w, nil)
		page.Execute(w, nil)
		FOOTER.Execute(w, nil)
	}
}
