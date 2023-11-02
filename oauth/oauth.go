package oauth

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

type OAuth[T any] interface {
	Login() (string, error)
	Callback(r *http.Request) (*oauth2.Token, error)
	GetUser(ctx context.Context, token string) (*T, error)
}

func New[T any](config oauth2.Config, userURL string, userAuthType TokenType) OAuth[T] {
	return &service[T]{
		config:       config,
		states:       NewGenerator(),
		userURL:      userURL,
		userAuthType: userAuthType,
	}
}
