package null_value

import (
	"database/sql"
	"time"
)

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}

	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NewStringFromNull(s sql.NullString) string {
	if s.Valid {
		return s.String
	}

	return ""
}

func NewNullDate(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

func NewTimeFromNull(t sql.NullTime) time.Time {
	if t.Valid {
		return t.Time
	}

	return time.Time{}
}

func NewNullInt(i int64) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{}
	}

	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

func NewIntFromNull(i sql.NullInt64) int64 {
	if i.Valid {
		return i.Int64
	}

	return 0
}

func NewNullFloat(f float64) sql.NullFloat64 {
	if f == 0 {
		return sql.NullFloat64{}
	}

	return sql.NullFloat64{
		Float64: f,
		Valid:   true,
	}
}
