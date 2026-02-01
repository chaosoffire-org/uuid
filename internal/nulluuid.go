package internal

import "database/sql/driver"

// NullUUID represents a UUID that may be null.
type NullUUID struct {
	UUID  UUID
	Valid bool
}

// Scan implements sql.Scanner.
func (nu *NullUUID) Scan(value interface{}) error {
	if value == nil {
		nu.UUID, nu.Valid = Nil, false
		return nil
	}

	nu.Valid = true

	return nu.UUID.Scan(value)
}

// Value implements driver.Valuer.
func (nu NullUUID) Value() (driver.Value, error) {
	if !nu.Valid {
		return nil, nil
	}

	return nu.UUID.Value()
}
