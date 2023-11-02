package sse

import "net/http"

type Upgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, event string)
}

type Publisher interface {
	Publish(event string, data []byte)
	TryPublish(event string, data []byte) bool

	PublishJSON(event string, data any) error
	TryPublishJSON(event string, data any) (bool, error)

	StreamExists(event string) bool
}

type Server interface {
	Upgrader
	Publisher
}
