package oauth

type TokenType int

const (
	AuthorizationHeader TokenType = iota
	AccessTokenQueryParam
)
