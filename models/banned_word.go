package models

import (
	"database/sql"
	"time"
)

// BannedWord represents a word or phrase that is banned in chat
type BannedWord struct {
	ID                   uint64 `gorm:"primaryKey"`
	OrganizationID       sql.NullInt64
	Organization         *Organization
	Word                 string
	TemporaryMuteSeconds sql.NullInt64
	PermanentBan         bool
	CreatedDate          time.Time
	DeletedDate          sql.NullTime
}
