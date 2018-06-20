package main

import (
	"net/http"
	"html/template"
)

func welcomeHandler(w http.ResponseWriter, r *http.Request){

	parsed_template, _ := template.ParseFiles("home.html")
	err := parsed_template.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}