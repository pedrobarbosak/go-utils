package models

type Discord struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	AvatarID string `json:"avatar"`
	Name     string `json:"global_name"`
	Email    string `json:"email"`
	MFA      bool   `json:"mfa_enabled"`
	Verified bool   `json:"verified"`
}
