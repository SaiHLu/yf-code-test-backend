package model

import (
	"time"
)

type UserLogEvent string

const (
	UserLogEventRead   UserLogEvent = "user:read"
	UserLogEventCreate UserLogEvent = "user:created"
	UserLogEventUpdate UserLogEvent = "user:updated"
	UserLogEventDelete UserLogEvent = "user:deleted"
)

func (e UserLogEvent) String() string {
	return string(e)
}

type UserLogModel struct {
	UserID    string       `json:"user_id" bson:"user_id"`
	Event     UserLogEvent `json:"event" bson:"event"`
	Data      interface{}  `json:"data,omitempty" bson:"data,omitempty"`
	CreatedAt time.Time    `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt *time.Time   `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt *time.Time   `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}
