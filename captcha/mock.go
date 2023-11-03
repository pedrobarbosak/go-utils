package captcha

import "context"

type mock struct{}

func (service *mock) Verify(string) error                             { return nil }
func (service *mock) VerifyWithContext(context.Context, string) error { return nil }
func NewMock() Service                                                { return &mock{} }
