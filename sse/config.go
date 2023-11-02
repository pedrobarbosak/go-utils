package sse

import (
	"time"

	"github.com/r3labs/sse/v2"
)

type Config struct {
	Headers      map[string]string
	EventTTL     time.Duration
	BufferSize   int
	EncodeBase64 bool
	SplitData    bool
	AutoStream   bool
	AutoReplay   bool
}

func NewConfig() Config {
	return Config{
		BufferSize: sse.DefaultBufferSize,
		AutoStream: true,
		AutoReplay: true,
		Headers:    map[string]string{},
	}
}

func newWithConfig(cfg Config) *sse.Server {
	s := sse.New()
	s.Headers = cfg.Headers
	s.EventTTL = cfg.EventTTL
	s.BufferSize = cfg.BufferSize
	s.EncodeBase64 = cfg.EncodeBase64
	s.SplitData = cfg.SplitData
	s.AutoStream = cfg.AutoStream
	s.AutoReplay = cfg.AutoReplay
	return s
}
