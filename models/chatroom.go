package models

import (
	"database/sql"
	"time"
)

// ChatRoom represents a single chat room, with a unique chat history
type ChatRoom struct {
	ID             uint64 `gorm:"primaryKey"`
	OrganizationID uint64
	Organization   *Organization
	Identifier     string
	Title          string
	CurrentUsers   int
	CreatedDate    time.Time
	DeletedDate    sql.NullTime
}
