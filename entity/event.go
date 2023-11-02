package entity

import "time"

type TimeEvent struct {
	Timestamp int64  `validate:"required"`
	UserID    string `validate:"required"`
}

func newTimeEvent(userID string) *TimeEvent {
	return &TimeEvent{
		Timestamp: time.Now().UTC().Unix(),
		UserID:    userID,
	}
}
