package main

import (
	"net/http"
	"html/template"
	"io/ioutil"
	"os"
	"fmt"
)

func downloadHandler(w http.ResponseWriter, req *http.Request) {

	server := getServerFromJSON()

	var i int = 0
	var scr string = ""
	var description string = ""
	var scriptsLocation string = server.DownloadLoc
	if _, err := os.Stat(scriptsLocation); os.IsExist(err) {
		os.Rename(scriptsLocation, "./seleniumScripts_backup/")
	}
	os.Mkdir(scriptsLocation, 0775)
	items := getItems()

	m := items.(map[string]interface{})
	for _, v := range m {
		switch vv := v.(type) {
		case string:
		case float64:
		case []interface{}:

			for _, u := range vv {

				a := u.(map[string]interface{})

				for l, m := range a {

					switch mm := m.(type){

					case string:
						if l == "description" {
							description += m.(string)
						}
					case []interface{}:
					//fmt.Println(mm)
					default:
						//fmt.Println(k,"is of a type I don't know how to handle")
						if mm !=  nil{}
						if l == "script" && m != nil {
							s := m.(map[string]interface{})
							for k, v := range s{
								if k == "script" {
									scr += v.(string)
								}
							}
						}
					}
				}
				//ioutil.WriteFile(scriptsLocation + "script" + strconv.Itoa(i) + ".py", []byte(scr), 0775)
				fmt.Println(description + ",    " + SanitizeFilename(description, false))
				ioutil.WriteFile(scriptsLocation + SanitizeFilename(description, false) + ".py", []byte(scr), 0775)
				description = ""
				scr = ""
				i += 1
			}
		//fmt.Println(scr)
		default:
		//fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
	parsed_template, _ := template.ParseFiles("downloadConfirmation.html")
	err := parsed_template.ExecuteTemplate(w, "downloadConfirmation.html", scriptsLocation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}