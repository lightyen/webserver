package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"webserver/server"
)

var (
	host string
	port string
)

func init() {
	flag.StringVar(&host, "host", "localhost", "listening host ip")
	flag.StringVar(&port, "p", "8080", "listening port number")
	flag.StringVar(&server.WebRoot, "w", "/www", "filepath to static web resouces")
	flag.BoolVar(&server.EnableGzip, "gzip", true, "enable gzip")
	flag.Parse()
}

func main() {
	srv := &http.Server{
		Addr:           net.JoinHostPort(host, port),
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

	fmt.Printf("\n >> %s http://localhost:%s/\n", "local:", port)

	if err := srv.Serve(ln); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
