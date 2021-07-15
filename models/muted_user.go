package models

import (
	"database/sql"
	"time"
)

// MutedUser is a user that is muted in chat
type MutedUser struct {
	ID             uint64 `gorm:"primaryKey"`
	OrganizationID uint64
	Organization   *Organization
	Username       string
	CreatedDate    time.Time
	DeletedDate    sql.NullTime
}
