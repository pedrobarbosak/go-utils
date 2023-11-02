package discord

import "time"

type Server struct {
	ID                       string   `json:"id"`
	Name                     string   `json:"name"`
	Icon                     string   `json:"icon"`
	Owner                    bool     `json:"owner"`
	Permissions              string   `json:"permissions"`
	Features                 []string `json:"features"`
	ApproximateMemberCount   int      `json:"approximate_member_count"`
	ApproximatePresenceCount int      `json:"approximate_presence_count"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	AvatarID string `json:"avatar"`
	Name     string `json:"global_name"`
	Email    string `json:"email"`
	MFA      bool   `json:"mfa_enabled"`
	Verified bool   `json:"verified"`
}

type ServerMember struct {
	User     User      `json:"user"`
	Nick     string    `json:"nick"`
	Avatar   *string   `json:"avatar"`
	Roles    []string  `json:"roles"`
	JoinedAt time.Time `json:"joined_at"`
	Deaf     bool      `json:"deaf"`
	Mute     bool      `json:"mute"`
}
