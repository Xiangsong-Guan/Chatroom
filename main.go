package main

import (
	"flag"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/Xiangsong-Guan/Chatroom/client"
	"github.com/Xiangsong-Guan/Chatroom/historylog"
	hs "github.com/Xiangsong-Guan/Chatroom/httpserver"
	hb "github.com/Xiangsong-Guan/Chatroom/hub"
	"github.com/Xiangsong-Guan/Chatroom/message"

	"github.com/gorilla/websocket"
)

// flags here
var addr = flag.String("addr", ":8080", "http service address")
var logFile = flag.String("log", "", "file path to log the chat history, leave it if just print to stdout")
var home = flag.String("home", "html/", "serve flie root path")
var indexFile = flag.String("index", "home.html", "index file in your html root dir, relative path")

// pre defined here
var spChar = message.SPChar

// some upvalue here
var hl *historylog.HistoryLog
var hub *hb.Hub

func main() {
	flag.Parse()

	*home = path.Clean(*home) + "/"
	log.Println("set html root:", *home)
	log.Println("set index file:", *home+*indexFile)

	// here we init the log
	hl = historylog.New(*logFile)
	if *logFile == "" {
		log.Println("fire up history log to stdout")
	} else {
		log.Println("fire up history log from:", *logFile)
	}
	defer hl.Close()

	// in hub implament: here we get hub pass hislog in
	hub = hb.New(hl)

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

	if r.URL.Path == "/" {
		http.ServeFile(w, r, *home+*indexFile)
	} else {
		http.ServeFile(w, r, *home+r.URL.Path)
	}

	log.Println("request:", r.Method, r.URL)
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
		http.Error(w, "Invalid personal infomation", http.StatusBadRequest)
		return
	}
	// here we tend to check this value in detail, but now, we just make a
	// simple nil check for tip.
	email, arimasu := r.Form["email"]
	if !arimasu {
		log.Println("invalid chatin form: email is nil")
		http.Error(w, "Invalid personal infomation", http.StatusBadRequest)
		return
	}
	nickname, arimasu := r.Form["nick"]
	if !arimasu {
		log.Println("invalid chatin form: nick is nil")
		http.Error(w, "Invalid personal infomation", http.StatusBadRequest)
		return
	}

	// here pass the check
	log.Println("welcome", email[0], "the", nickname[0])

	/* we get this kind of personal info keep by user's cookie.
	 * in next chat page, user's client fire up a ws conn,
	 * first it shuold pass this kind of info again, and we new a
	 * client at that time. */

	//here we server chat page
	http.ServeFile(w, r, *home+"chatroom.html")
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

	// in client implament: make new client here
	_, m, err := conn.ReadMessage()
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("fail to assgin resource"))
		log.Println("fail to initial remote info:", err)
		return
	}
	info := strings.Split(string(m), message.SPChar)
	cInfo := message.Info{
		Name:      info[0],
		Mail:      info[1],
		IPAddress: strings.Split(conn.RemoteAddr().String(), ":")[0],
	}
	c := client.New(hub, conn, cInfo)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go c.WritePump()
	go c.ReadPump()
	log.Println(c.GenSummary())
}
