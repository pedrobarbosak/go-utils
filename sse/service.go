package sse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/r3labs/sse/v2"
)

type server struct {
	server *sse.Server
}

func (s *server) Upgrade(w http.ResponseWriter, r *http.Request, event string) {
	const stream = "stream="

	if !strings.Contains(r.URL.RawQuery, stream+event) {
		r.URL.RawQuery = fmt.Sprintf("stream=%s&%s", event, r.URL.RawQuery)
	}

	s.server.ServeHTTP(w, r)
}

func (s *server) Publish(event string, data []byte) {
	s.server.Publish(event, &sse.Event{Data: data})
}

func (s *server) TryPublish(event string, data []byte) bool {
	return s.server.TryPublish(event, &sse.Event{Data: data})
}

func (s *server) PublishJSON(event string, data any) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	s.server.Publish(event, &sse.Event{Data: d})

	return nil
}

func (s *server) TryPublishJSON(event string, data any) (bool, error) {
	d, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	return s.server.TryPublish(event, &sse.Event{Data: d}), nil
}

func (s *server) StreamExists(event string) bool {
	return s.server.StreamExists(event)
}

func New(cfg ...Config) Server {
	if len(cfg) == 0 {
		return &server{
			server: sse.New(),
		}
	}

	return &server{
		server: newWithConfig(cfg[0]),
	}
}

func NewWithCallback(onSubscribe, onUnsubscribe func(streamID string, sub *sse.Subscriber), cfg ...Config) Server {
	if len(cfg) == 0 {
		return &server{
			server: sse.NewWithCallback(onSubscribe, onUnsubscribe),
		}
	}

	s := newWithConfig(cfg[0])
	s.OnSubscribe = onSubscribe
	s.OnUnsubscribe = onUnsubscribe

	return &server{
		server: s,
	}
}
