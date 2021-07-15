package models

import (
	"database/sql"
	"time"
)

// Organization is a company or individual profile, which can contain
// multiple chat rooms within it.
type Organization struct {
	ID          uint64 `gorm:"primaryKey"`
	AccountID   uint64
	Account     *Account
	Name        string
	CreatedDate time.Time
	DeletedDate sql.NullTime
}
