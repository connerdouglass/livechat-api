package utils

import "database/sql"

// FlattenNullString converts a nullstring to either nil, or a string pointer
func FlattenNullString(str sql.NullString) *string {
	if !str.Valid {
		return nil
	}
	return &(str.String)
}

// FlattenNullInt64 converts a nullable int64 to either nil, or a pointer to the int64
func FlattenNullInt64(val sql.NullInt64) *int64 {
	if !val.Valid {
		return nil
	}
	return &(val.Int64)
}

// FlattenNullInt32 converts a nullable int64 to either nil, or a pointer to the int32
func FlattenNullInt32(val sql.NullInt32) *int32 {
	if !val.Valid {
		return nil
	}
	return &(val.Int32)
}

// FlattenNullTimeMilli converts a nullable timestamp to either null, or a pointer to the time in milliseconds
func FlattenNullTimeMilli(val sql.NullTime) *uint64 {
	if !val.Valid {
		return nil
	}
	ms := uint64(val.Time.UTC().Unix() * 1000)
	return &ms
}

// FlattenNullTimeSec converts a nullable timestamp to either null, or a pointer to the time in seconds
func FlattenNullTimeSec(val sql.NullTime) *uint64 {
	if !val.Valid {
		return nil
	}
	ms := uint64(val.Time.UTC().Unix())
	return &ms
}
