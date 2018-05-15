package main

import (
	hs "chatroom/httpserver"
	"flag"
	"log"
	"net/http"
)

// flags here
var addr = flag.String("addr", ":8080", "http service address")
var logFile = flag.String("log file", "", "file path to log the chat history, leave it if just print to stdout")

func main() {
	flag.Parse()

	// here we init the log
	historyLog

	// here we get hub

	// here we add handler to its url
	var handlers = map[string]func(http.ResponseWriter, *http.Request){
		"/": serveHome,
	}

	hs.New(addr, handlers).Run()
}

// here we will make some handlers
// serveHome handle "/", it just server static file home.html
func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		log.Println("invalid url:", r.URL.Path, "for \"/\"")
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		log.Println("invalid method:", r.Method, "for \"GET\"")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
	log.Println("servered", r.Method, r.URL)
}
