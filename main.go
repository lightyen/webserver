package main

import (
	"flag"
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
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
