package main

import (
	"Chatroom/historylog"
	hs "Chatroom/httpserver"
	hb "Chatroom/hub"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
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

	// TODO in hub implament: here we get hub pass hislog in

	// here we add handler to its url
	var handlers = map[string]func(http.ResponseWriter, *http.Request){
		"/":           serveStatic,
		"/chatin":     chatin,
		"/chatsocket": serveWs,
	}

	log.Println("get httpserver listen on:", *addr)
	hs.New(addr, handlers).Run()
}

// here we will make some handlers

// serveStatic handle static file.
// if get "/", it just server static file home.html
func serveStatic(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Println("invalid method:", r.Method, "for static \"GET\"")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("request:", r.Method, r.URL)

	if r.URL.Path == "/" {
		http.ServeFile(w, r, "html/home.html")
	} else {
		http.ServeFile(w, r, "html"+r.URL.Path)
	}
}

// chatin handle post method request to new a client
func chatin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Println("invalid method:", r.Method, "for chatin")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("request:", r.Method, r.URL)

	// parse and check post
	if err := r.ParseForm(); err != nil {
		log.Println("fail to parse this req's form:", err)
	}
	// here we tend to check this value in detail, but now, we just make a
	// simple nil check for tip.
	email, arimasu := r.Form["email"]
	if !arimasu {
		log.Println("invalid chatin form: email is nil")
		http.Error(w, "Invalid personal infomation", http.StatusBadRequest)
	}
	nickname, arimasu := r.Form["nick"]
	if !arimasu {
		log.Println("invalid chatin form: nick is nil")
		http.Error(w, "Invalid personal infomation", http.StatusBadRequest)
	}
	log.Println("welcome", email, "the", nickname)

	/* we get this kind of personal info keep by user's cookie.
	 * in next chat page, user's client fire up a ws conn,
	 * first it shuold pass this kind of info again, and we new a
	 * client at that time. */

	// TODO in ws handle: here we server chat page
}

// serveWs handle websocket conn, make a client and let this client
// take over.
func serveWs(w http.ResponseWriter, r *http.Request) {
	log.Println("request websocket")

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("fail to make ws:", err)
		return
	}

	// TODO in client implament: make new client here
	log.Println(conn)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	// go client.WritePump()
	// go client.ReadPump()
}
