package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	// path is relative to current working directory
	files := []string{
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/base.tmpl.html",
		"./ui/html/nav.tmpl.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	rawId := r.PathValue("id")
	id, err := strconv.Atoi(rawId)

	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with the ID %d", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a snippet"))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet"))
}
