package main

import (
	"net/http"
	"html/template"
	"strings"
	"strconv"
)

func jobDetailHandler(w http.ResponseWriter, req *http.Request) {

	scriptId := strings.Replace(req.URL.RequestURI(), "/jobDetail/", "", -1)

	items1 := getJob(scriptId)
	job1 := make(map[string]string,50)
	myJob1 := make(map[string]string,50)

	elements := items1.(map[string]interface{})
	for l, m := range elements {
		switch mm := m.(type) {
		case string:
			job1[l] = m.(string)
		case bool:
			job1[l] = strconv.FormatBool(m.(bool))
		case float64:
			job1[l] = strconv.FormatFloat(m.(float64), 'f', 0, 64)
		case []interface{}:
			var temp string = "["
			for _, e := range mm {
				temp += " " + e.(string)
			}
			temp += " ]"
			job1[l] = temp
		default:
			if mm != nil{}
			if m != nil {
				s := m.(map[string]interface{})
				for k, v := range s{
					switch oo := v.(type){
					case string:
						if l != "script" && k != "script" {
							job1[l+"--"+k] = v.(string)
						}
					case float64:
						job1[l+"--"+k] = strconv.FormatFloat(v.(float64), 'f', 0, 64)
					default:
						if oo != nil{}
					}

				}
			}

		}
	}
	if job1["_id"] == scriptId {
		for k, v := range job1 {
			myJob1[k] = v
		}
	}

	parsed_template, _ := template.ParseFiles("jobDetail.html")
	err := parsed_template.ExecuteTemplate(w, "jobDetail.html", myJob1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}