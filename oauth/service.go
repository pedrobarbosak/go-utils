package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

type service[T any] struct {
	config       oauth2.Config
	states       StateGenerator
	userURL      string
	userAuthType TokenType
}

func (s *service[T]) Login() (string, error) {
	state, err := s.states.New()
	if err != nil {
		return "", err
	}

	return s.config.AuthCodeURL(state), nil
}

func (s *service[T]) Callback(r *http.Request) (*oauth2.Token, error) {
	state := r.FormValue("state")
	if err := s.states.Validate(state); err != nil {
		return nil, err
	}

	code := r.FormValue("code")
	token, err := s.config.Exchange(r.Context(), code)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *service[T]) GetUser(ctx context.Context, token string) (*T, error) {
	req, err := http.NewRequest("GET", s.userURL, nil)
	if err != nil {
		return nil, err
	}

	s.prepareAuthenticationRequest(req, token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Println(string(data))

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid response")
	}

	var obj T
	if err = json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	return &obj, nil
}

func (s *service[T]) prepareAuthenticationRequest(req *http.Request, token string) {
	switch s.userAuthType {

	case AuthorizationHeader:
		req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	case AccessTokenQueryParam:
		q := req.URL.Query()
		q.Add("access_token", token)
		req.URL.RawQuery = q.Encode()
	}
}
