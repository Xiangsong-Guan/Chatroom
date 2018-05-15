package main

import (
	"Chatroom/historylog"
	hs "Chatroom/httpserver"
	hb "Chatroom/hub"
	"flag"
	"log"
	"net/http"
)

// flags here
var addr = flag.String("addr", ":8080", "http service address")
var logFile = flag.String("log file", "", "file path to log the chat history, leave it if just print to stdout")

// some upvalue here
var hl historylog.HistoryLog
var hub hb.Hub

func main() {
	log.Println("Parse flags...")
	flag.Parse()

	// here we init the log
	hl = historylog.New(*logFile)
	log.Println("fire up history log from:", *logFile)
	defer hl.Close()

	// here we get hub

	// here we add handler to its url
	var handlers = map[string]func(http.ResponseWriter, *http.Request){
		"/": serveHome,
	}

	log.Println("get httpserver listen on:", *addr)
	hs.New(addr, handlers).Run()
}

// here we will make some handlers
// serveHome handle "/", it just server static file home.html
func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Println("invalid method:", r.Method, "for \"GET\"")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	switch r.URL.Path {
	case "/":
		http.ServeFile(w, r, "html/home.html")
	case "/css/home.css":
		http.ServeFile(w, r, "html/css/home.css")
	case "/js/home.js":
		http.ServeFile(w, r, "html/js/home.js")
	default:
		log.Println("invalid url:", r.URL)
	}

	log.Println("servered", r.Method, r.URL)
}
