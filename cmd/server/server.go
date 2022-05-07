package main

import (
	"flag"
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/gerifield/obs-overlay/server"
)

func main() {
	listen := flag.String("listen", ":8080", "Server listen address")
	static := flag.String("static", "./static", "Static folder")
	flag.Parse()

	conf := server.Config{
		StaticDir: *static,
	}

	logger, err := zap.NewProduction(zap.AddStacktrace(zap.FatalLevel))
	if err != nil {
		log.Println(err)
		return
	}

	defer logger.Sync()

	srv := server.New(logger, conf)

	log.Println("listening on", *listen)
	http.ListenAndServe(*listen, srv.Routes())
}
