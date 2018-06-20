package main

import (
	"net/http"
	"html/template"
)

func listHandler(w http.ResponseWriter, r *http.Request){
	items := getSchedules()

	parsed_template, _ := template.ParseFiles("listofjobs.html")
	err := parsed_template.ExecuteTemplate(w, "listofjobs.html", items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}