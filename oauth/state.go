package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

type StateGenerator interface {
	New() (string, error)
	Validate(state string) error
}

type generator struct {
	mutex   *sync.Mutex
	states  map[string]time.Time
	timeout time.Duration
}

func NewGenerator() StateGenerator {
	return &generator{
		mutex:   &sync.Mutex{},
		states:  make(map[string]time.Time),
		timeout: time.Minute * 3,
	}
}

func (s *generator) New() (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var err error
	for i := 0; i < 10; i++ {
		var state string
		state, err = s.new()
		if err != nil {
			continue
		}

		s.states[state] = time.Now()
		return state, nil
	}

	return "", err
}

func (s *generator) Validate(state string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	t, ok := s.states[state]
	if !ok {
		return errors.New("invalid state")
	}

	delete(s.states, state)

	if time.Now().After(t.Add(s.timeout)) {
		return errors.New("timout state")
	}

	return nil
}

func (s *generator) new() (string, error) {
	b := make([]byte, 256)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
