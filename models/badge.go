package models

import (
	"database/sql"
	"time"
)

// Badge is a symbol that can be assigned to users
type Badge struct {
	ID              uint64 `gorm:"primaryKey"`
	OrganizationID  uint64
	Organization    *Organization
	BackgroundColor string
	Image           string
	CreatedDate     time.Time
	DeletedDate     sql.NullTime
}
