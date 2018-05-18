// Package httpserver basic http server for chatroom
package httpserver

import (
	"log"
	"net/http"
)

// HTTPServer struct
type HTTPServer struct {
	listenAddr string
	handlers   map[string]func(http.ResponseWriter, *http.Request)
}

// New generate a HTTPServer obj with listen ipaddress and some handlers
// handler are map, its key(string) reffer to url
func New(addr *string, handlers map[string]func(http.ResponseWriter, *http.Request)) HTTPServer {
	return HTTPServer{
		listenAddr: *addr,
		handlers:   handlers,
	}
}

// Run is a great elder who can make everything simple
func (hs HTTPServer) Run() {
	// make sure we can handle something
	if _, arimasu := hs.handlers["/"]; !arimasu {
		log.Println("No index handler", hs)
	}

	// register handler to pattern
	for pattern, handler := range hs.handlers {
		http.HandleFunc(pattern, handler)
	}

	// run http server
	if err := http.ListenAndServe(hs.listenAddr, nil); err != nil {
		log.Fatalln("Server fail to start", err)
	}
}

// ListenAddr server's listening address
func (hs HTTPServer) ListenAddr() string {
	return hs.listenAddr
}
