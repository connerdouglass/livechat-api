package main

import (
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ParseDatabaseDriver parses a database string to return the appropriate driver
func ParseDatabaseDriver(dsn string) gorm.Dialector {

	// If it's SQLite protocol
	if strings.HasPrefix(dsn, "sqlite://") {

		// Parse out the filename
		filename := strings.TrimPrefix(dsn, "sqlite://")

		// Open the connection to the file
		return sqlite.Open(filename)

	}

	// If it's MySQL protocol
	if strings.HasPrefix(dsn, "mysql://") {

		// Parse out the database url
		dbUrl := strings.TrimPrefix(dsn, "mysql://")

		// Return the connection
		return mysql.Open(dbUrl)

	}

	// Return nil otherwise
	return nil

}
