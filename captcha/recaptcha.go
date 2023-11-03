package captcha

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Recaptcha struct {
	url       string
	secretKey string
	threshold float32
}

type Result struct {
	Success    bool     `json:"success"`
	Timestamp  string   `json:"challenge_ts"`
	Hostname   string   `json:"hostname"`
	ErrorCodes []string `json:"error-codes"`
	Score      float32  `json:"score"`
}

func NewRecaptcha(secretKey string, threshold float32) Service {
	return &Recaptcha{
		url:       "https://www.google.com/recaptcha/api/siteverify",
		secretKey: secretKey,
		threshold: threshold,
	}
}

func (s *Recaptcha) Verify(response string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.post(ctx, response)
}

func (s *Recaptcha) VerifyWithContext(ctx context.Context, response string) error {
	return s.post(ctx, response)
}

func (s *Recaptcha) post(ctx context.Context, response string) error {
	data := url.Values{"secret": {s.secretKey}, "response": {response}}

	req, err := http.NewRequest("POST", s.url, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}

	var result Result
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if !result.Success {
		for _, code := range result.ErrorCodes {
			err = errors.Join(err, errors.New(code))
		}
		return err
	}

	if result.Score < s.threshold {
		return fmt.Errorf("score below threshold: %.1f", result.Score)
	}

	return nil
}
