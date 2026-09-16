package database

import "github.com/oklog/ulid/v2"

// GenerateULID returns a new ULID as a string.
func GenerateULID() string {
	return ulid.Make().String()
}
