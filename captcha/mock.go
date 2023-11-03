package captcha

type mock struct{}

func (service *mock) Verify(string) error { return nil }
func NewMock() Service                    { return &mock{} }
