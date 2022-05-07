package server

import (
	"context"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"

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

	connectionsLock *sync.Mutex
	connections     map[int64]chan string
}

func New(logger *zap.Logger, conf Config) *Logic {
	l := &Logic{
		config:          conf,
		logger:          logger,
		event:           make(chan string),
		connectionsLock: &sync.Mutex{},
		connections:     make(map[int64]chan string),
	}

	// go func() {
	// 	ticker := time.NewTicker(1 * time.Second)
	// 	i := 0
	// 	for range ticker.C {
	// 		l.event <- fmt.Sprintf("hello-%d", i)
	// 		i++
	// 	}
	// }()
	go l.eventLoop()

	return l
}

func (l *Logic) eventLoop() {
	for evt := range l.event {
		l.connectionsLock.Lock()
		for id, ch := range l.connections {
			select {
			case ch <- evt:
			default:
				l.logger.Warn("message skipped", zap.Int64("cid", id))
			}

		}
		l.connectionsLock.Unlock()
	}
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

	connectionID := rand.Int63n(1000000000)
	logger := l.logger.With(zap.Int64("cid", connectionID))
	logger.Info("new connection")

	event := make(chan string, 20)
	l.registerConnection(connectionID, event)
	l.logStats()

	defer func() {
		l.unregisterConnection(connectionID)
		close(event)
		l.logStats()
	}()

	go func() {
		for evt := range event {
			logger.Info("send event", zap.String("event", evt))
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err = c.Write(ctx, websocket.MessageText, []byte(evt))
			if err != nil {
				logger.Error("websocket write failed", zap.Error(err))
				cancel()

				return
			}
			cancel()
		}
	}()

	for {
		_, _, err := c.Read(context.Background())
		if err != nil {
			logger.Error("websocket read failed", zap.Error(err))
			return
		}
	}

	// c.Close(websocket.StatusNormalClosure, "Bye")
}

func (l *Logic) logStats() {
	l.connectionsLock.Lock()
	conns := len(l.connections)
	l.connectionsLock.Unlock()

	l.logger.Info("stats", zap.Int("connectionNum", conns))
}

func (l *Logic) registerConnection(cid int64, ch chan string) {
	l.connectionsLock.Lock()
	l.connections[cid] = ch
	l.connectionsLock.Unlock()
}

func (l *Logic) unregisterConnection(cid int64) {
	l.connectionsLock.Lock()
	delete(l.connections, cid)
	l.connectionsLock.Unlock()
}
