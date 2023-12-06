package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"webserver/server"
)

var port string

func init() {
	flag.StringVar(&port, "p", "8080", "listening port number")
	flag.StringVar(&server.WebRoot, "w", "www", "filepath to static web resouces")
	flag.BoolVar(&server.EnableGzip, "gzip", true, "enable gzip")
	flag.Parse()
}

func main() {
	srv := &http.Server{
		Addr:           net.JoinHostPort("", port),
		Handler:        server.NewRouter(),
		MaxHeaderBytes: 1 << 20,
	}
	if srv.Addr == "" {
		srv.Addr = ":http"
	}

	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("listening on:", srv.Addr)

	if err := srv.Serve(ln); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
