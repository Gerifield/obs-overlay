package server

import (
	"context"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
	"nhooyr.io/websocket"
)

type Config struct {
	StaticDir string
}

type Logic struct {
	config Config
	logger *zap.Logger

	event chan string
}

func New(logger *zap.Logger, conf Config) *Logic {
	l := &Logic{
		config: conf,
		logger: logger,
		event:  make(chan string),
	}

	// go func() {
	// 	ticker := time.NewTicker(1 * time.Second)
	// 	i := 0
	// 	for range ticker.C {
	// 		l.event <- fmt.Sprintf("hello-%d", i)
	// 		i++
	// 	}
	// }()

	return l
}

func (l *Logic) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(l.config.StaticDir))

	mux.Handle("/static/", http.StripPrefix("/static", fs))
	mux.HandleFunc("/websocket", l.wsHandler)
	mux.HandleFunc("/http", l.httpHandler)

	return mux
}

func (l *Logic) Stop() {
	close(l.event)
}

func (l *Logic) httpHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		l.logger.Error("read HTTP failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	l.event <- string(b)
}

func (l *Logic) wsHandler(w http.ResponseWriter, r *http.Request) {
	l.logger.Info("ws connection")

	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		l.logger.Error("websocket accept failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	for evt := range l.event {
		err = c.Write(context.Background(), websocket.MessageText, []byte(evt))
		if err != nil {
			l.logger.Error("websocket write failed", zap.Error(err))

			return
		}
	}

	c.Close(websocket.StatusNormalClosure, "Bye")
}
