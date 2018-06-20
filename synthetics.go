package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"fmt"
	"net/http"
	"time"
)

var serviceUri string = ""

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

func getJob(id string) interface{} {
	var items interface{}
	eum := getEumFromJSON()
	req, _ := http.NewRequest("GET", serviceUri + id, nil)
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

