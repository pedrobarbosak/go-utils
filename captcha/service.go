package captcha

import "context"

type Service interface {
	Verify(response string) error
	VerifyWithContext(ctx context.Context, response string) error
}
