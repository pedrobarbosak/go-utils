package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetUserServes(ctx context.Context, token string) ([]*Server, error) {
	url := "https://discord.com/api/v10/users/@me/guilds"

	body, err := makeRequest(ctx, url, token)
	if err != nil {
		return nil, err
	}

	servers := make([]*Server, 0)
	if err = json.Unmarshal(body, &servers); err != nil {
		return nil, err
	}

	return servers, nil
}

func GetUserServerProfile(ctx context.Context, serverID string, token string) (*ServerMember, error) {
	url := fmt.Sprintf("https://discord.com/api/v10/users/@me/guilds/%s/member", serverID)

	body, err := makeRequest(ctx, url, token)
	if err != nil {
		return nil, err
	}

	member := &ServerMember{}
	if err = json.Unmarshal(body, member); err != nil {
		return nil, err
	}

	return member, nil
}

func makeRequest(ctx context.Context, url string, token string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}
