package oauth

import (
	"github.com/pedrobarbosak/go-utils/oauth/models"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

func NewGoogle(clientID string, secret string, redirectURL string, scopes ...string) OAuth[models.Google] {
	return New[models.Google](oauth2.Config{
		ClientID:     clientID,
		ClientSecret: secret,
		Endpoint:     endpoints.Google,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
	}, "https://www.googleapis.com/oauth2/v2/userinfo",
		AccessTokenQueryParam,
	)
}

func NewFacebook(clientID string, secret string, redirectURL string, scopes ...string) OAuth[models.Facebook] {
	return New[models.Facebook](oauth2.Config{
		ClientID:     clientID,
		ClientSecret: secret,
		Endpoint:     endpoints.Facebook,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
	}, "https://graph.facebook.com/me?fields=id,name,email,picture",
		AccessTokenQueryParam,
	)
}

func NewMicrosoft(clientID string, secret string, redirectURL string, tenant string, scopes ...string) OAuth[models.Microsoft] {
	return New[models.Microsoft](oauth2.Config{
		ClientID:     clientID,
		ClientSecret: secret,
		Endpoint:     endpoints.AzureAD(tenant),
		RedirectURL:  redirectURL,
		Scopes:       scopes,
	}, "https://graph.microsoft.com/v1.0/me",
		AuthorizationHeader,
	)
}

func NewDiscord(clientID string, secret string, redirectURL string, scopes ...string) OAuth[models.Discord] {
	return New[models.Discord](oauth2.Config{
		ClientID:     clientID,
		ClientSecret: secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
		RedirectURL: redirectURL,
		Scopes:      scopes,
	}, "https://discord.com/api/v10/users/@me",
		AuthorizationHeader,
	)
}
