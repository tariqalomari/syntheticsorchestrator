package main

import (
	"net/http"
	"strings"
	"fmt"
)

func getScriptHandler(w http.ResponseWriter, req *http.Request) {

	scriptId := strings.Replace(req.URL.RequestURI(), "/script/", "", -1)
	var scr string = ""
	var i int = 0
	items := getJob(scriptId)

				a := items.(map[string]interface{})

				for l, m := range a {

					switch mm := m.(type){

					case string:
						if l == "_id" {
							if m == scriptId {
								i += 1
							}
						}
					case []interface{}:
					//fmt.Println(mm)
					default:
						//fmt.Println(k,"is of a type I don't know how to handle")
						if mm !=  nil{}
						if l == "script" {
							//fmt.Println(l, m)
							s := m.(map[string]interface{})
							for k, v := range s{
								if k == "script" {
									scr += v.(string)
									i = 0
								}
							}
						}
					}
				}
			
	var filename string = "attachment; filename=" + scriptId + ".py"
	w.Header().Set("Content-Disposition", filename)
	w.Header().Set("Content-Type", "application/python")
	fmt.Fprintf(w, "%s", scr)

}
