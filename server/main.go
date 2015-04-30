package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Config struct {
	CurrentVersionString string
	TeamCodes            []string
	Address              string
}

var (
	config *Config
	cLock  sync.Mutex
)

func LoadConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(file, &config)
	return &config, err
}

func main() {
	fmt.Println("Starting abpatcher server")
	c, err := LoadConfig("config.json")
	if err != nil {
		panic(err)
	}
	config = c

	go ConfigRefresher("config.json")

	fmt.Println("Starting http server")
	http.HandleFunc("/version", HandleGetVersionString)
	err = http.ListenAndServe(config.Address, nil)
	if err != nil {
		panic(err)
	}
}

// Refreshes the config every 5 seconds
func ConfigRefresher(path string) {
	ticker := time.NewTicker(5 * time.Second)
	for {
		cLock.Lock()
		c, err := LoadConfig(path)
		if err != nil {
			fmt.Println("Error loading config: ", err)
			cLock.Unlock()
			return
		}
		config = c
		cLock.Unlock()
		<-ticker.C
	}
}

type Response struct {
	Error   string
	Version string
}

func Respond(err, version string, statusCode int, w http.ResponseWriter, r *http.Request) {
	reponse := Response{
		Error:   err,
		Version: version,
	}

	serialized, rerr := json.Marshal(&reponse)
	if rerr != nil {
		fmt.Println("Error serializing response: ", rerr)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(serialized)
}

func HandleGetVersionString(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	teamCode := params.Get("tc")
	// Check if it matches up
	cLock.Lock()
	defer cLock.Unlock()
	found := false
	for _, v := range config.TeamCodes {
		if v == teamCode {
			found = true
			break
		}
	}

	if !found {
		fmt.Println("Someone tried to use incorrect teamcode: ", teamCode)
		Respond("Incorrect teamcode", "", http.StatusBadRequest, w, r)
		return
	}

	Respond("", config.CurrentVersionString, http.StatusOK, w, r)
}
