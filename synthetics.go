package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"fmt"
	"net/http"
	"html/template"
	"time"
	"strings"
	"strconv"
)

var serviceUri string = ""

func getScriptHandler(w http.ResponseWriter, req *http.Request) {

	scriptId := strings.Replace(req.URL.RequestURI(), "/script/", "", -1)
	var scr string = ""
	var i int = 0
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

						if i == 1 {
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
				}

			}
			//fmt.Println(scr)
		default:
			//fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
	var filename string = "attachment; filename=" + scriptId + ".py"
	w.Header().Set("Content-Disposition", filename)
	w.Header().Set("Content-Type", "application/python")
	fmt.Fprintf(w, "%s", scr)

}

func downloadHandler2(w http.ResponseWriter, req *http.Request) {

}

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

func jobDetailHandler(w http.ResponseWriter, req *http.Request) {

	scriptId := strings.Replace(req.URL.RequestURI(), "/jobDetail/", "", -1)
	job := make(map[string]string,50)
	myJob := make(map[string]string,50)
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
						job[l] = m.(string)

					case bool:
						job[l] = strconv.FormatBool(m.(bool))

					case float64:
						job[l] = strconv.FormatFloat(m.(float64), 'f', 0, 64)
					case []interface{}:
						//fmt.Printf("got %T\n", mm)
						var temp string = "["
						for _, e := range mm {
							temp += " " + e.(string)
						}
						temp += " ]"
						job[l] = temp
					default:
						//fmt.Println(k,"is of a type I don't know how to handle")
						if mm != nil{}
						if m != nil {
							s := m.(map[string]interface{})
							for k, v := range s{
								switch oo := v.(type){
								case string:
									if l != "script" && k != "script" {
										job[l+"--"+k] = v.(string)
									}
								case float64:
									job[l+"--"+k] = strconv.FormatFloat(v.(float64), 'f', 0, 64)
								default:
									if oo != nil{}
								}
							}
						}

					}
				}
				if job["_id"] == scriptId {
					for k, v := range job {
						myJob[k] = v
					}
				}

			}
		//fmt.Println(scr)
		default:
		//fmt.Println(k, "is of a type I don't know how to handle")
		}
	}

	parsed_template, _ := template.ParseFiles("jobDetail.html")
	err := parsed_template.ExecuteTemplate(w, "jobDetail.html", myJob)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func listHandler(w http.ResponseWriter, r *http.Request){
	items := getSchedules()

	parsed_template, _ := template.ParseFiles("listofjobs.html")
	err := parsed_template.ExecuteTemplate(w, "listofjobs.html", items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func welcomeHandler(w http.ResponseWriter, r *http.Request){

	parsed_template, _ := template.ParseFiles("home.html")
	err := parsed_template.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

type DataModel struct {
	Title string
}

func gui_data() *DataModel {
	gui_data := &DataModel{
		Title:	"My Title",
	}

	return gui_data
}

func main() {

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/script/", getScriptHandler)
	http.HandleFunc("/jobDetail/", jobDetailHandler)
	http.HandleFunc("/downloadScripts", downloadHandler)

	server := getServerFromJSON()

	serviceUri = server.ServiceUri

	s := &http.Server{
		Addr:		server.Port,
		ReadTimeout:	10 * time.Second,
		WriteTimeout:	10 * time.Second,
		MaxHeaderBytes:	1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}

// APIError to get HTTP response code to expected errors
type APIError struct {
	Message string
	Code    int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Message)
}

type Items struct {
	Items []Item	`json:"_items"`
}

type scriptItem struct {
	contentType string `json:"contentType"`
	seleniumScript string `json:"script"`
}

type Item struct {
	Id string `json:"_id"`
	UserEnabled bool `json:"userEnabled"`
	SystemEnabled bool `json:"systemEnabled"`
	TimeoutSeconds int `json:"timeoutSeconds"`
	Name string `json:"description"`

	scriptitem scriptItem `json:"script"`
}


func getSchedules() Items {
	eum := getEumFromJSON()
	//ctx := context.WithTimeout(context.Background(), 5 * time.Second)
	//req = req.WithContext(ctx)
	req, _ := http.NewRequest("GET", serviceUri, nil)
	req.SetBasicAuth(eum.Username, eum.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		err := &APIError{
			Code:    resp.StatusCode,
			Message: fmt.Sprintf("Status Code Error: %d\nRequest: %v", resp.StatusCode, req),
		}
		log.Printf("Status code greater than 400")
		log.Printf(err.Message)
		//return err
	}

	htmlData, err := ioutil.ReadAll(resp.Body) //<--- here!

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var items Items

	json.Unmarshal(htmlData, &items)

	return items

}

func getItems() interface{}{

	var items interface{}
	eum := getEumFromJSON()
	req, _ := http.NewRequest("GET", serviceUri, nil)
	req.SetBasicAuth(eum.Username, eum.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := http.DefaultClient.Do(req)

	if resp == nil {
		log.Printf("Response object is Nil")
		os.Exit(1)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		err := &APIError {
			Code:	resp.StatusCode,
			Message:	fmt.Sprintf("Status Code Error: %d\nRequest: %v", resp.StatusCode, req),
		}
		log.Printf(err.Message)
	}

	htmlData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err1 := json.Unmarshal(htmlData, &items)

	if err1 != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return items

}

// HTTP Server information
type ServerConfig struct {
	Port		string `json:"port"`
	DownloadLoc     string `json:"scriptsDownloadLocation"`
	ServiceUri 	string `json:"url"`
}

// create Server object from the serverconf.json file
func getServerFromJSON() ServerConfig {
	absPath, _ := filepath.Abs("serverconf.json")
	raw, err := ioutil.ReadFile(absPath)
	if err != nil {
		panic(err.Error())
	}

	var srv ServerConfig
	err = json.Unmarshal(raw, &srv)
	if err != nil {
		panic(err.Error())
	}
	return srv
}

// EumConfig is the Eum account where we want to send the data to
type EumConfig struct {
	Username              string `json:"username"`
	Password              string `json:"password"`
	Url                   string `json:"url"`
	AuthorizationHeader   string `json:"authorizationHeader"`
}

//var client = appdrest.NewClient("http", "demo2.appdynamics.com", 80, "demouser", "Ghed7ped0geN", "customer1")

// create EUM object from the eumconf.json file
func getEumFromJSON() EumConfig {
	absPath, _ := filepath.Abs("eumconf.json")
	raw, err := ioutil.ReadFile(absPath)
	if err != nil {
		panic(err.Error())
	}

	var eum EumConfig
	err = json.Unmarshal(raw, &eum)
	if err != nil {
		panic(err.Error())
	}
	return eum
}

var badCharacters = []string{
	"../",
	"https://",
	"http://",
	"<!--",
	"-->",
	"<",
	">",
	"'",
	"\"",
	"&",
	"$",
	"#",
	"{", "}", "[", "]", "=",
	";", "?", "%20", "%22",
	"%3c",   // <
	"%253c", // <
	"%3e",   // >
	"",   // > -- fill in with % 0 e - without spaces in between
	"%28",   // (
	"%29",   // )
	"%2528", // (
	"%26",   // &
	"%24",   // $
	"%3f",   // ?
	"%3b",   // ;
	"%3d",   // =
}

func RemoveBadCharacters(input string, dictionary []string) string {

	temp := input

	for _, badChar := range dictionary {
		temp = strings.Replace(temp, badChar, "", -1)
	}
	return temp
}

func SanitizeFilename(name string, relativePath bool) string {

	// default settings
	var badDictionary []string = badCharacters

	if name == "" {
		return name
	}

	// if relativePath is TRUE, we preserve the path in the filename
	// If FALSE and will cause upper path foldername to merge with filename
	// USE WITH CARE!!!

	if !relativePath {
		// add additional bad characters
		badDictionary = append(badCharacters, "./")
		badDictionary = append(badDictionary, "/")
	}

	// trim(remove)white space
	trimmed := strings.TrimSpace(name)

	// trim(remove) white space in between characters
	trimmed = strings.Replace(trimmed, " ", "", -1)

	// remove bad characters from filename
	trimmed = RemoveBadCharacters(trimmed, badDictionary)

	stripped := strings.Replace(trimmed, "\\", "", -1)

	return stripped
}