package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	listen := flag.String("listen", ":8080", "Server listen address")
	static := flag.String("static", "./static", "Static folder")
	flag.Parse()

	fs := http.FileServer(http.Dir(*static))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	log.Println("listening on", *listen)
	http.ListenAndServe(*listen, nil)
}
