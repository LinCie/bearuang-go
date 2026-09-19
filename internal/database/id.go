package database

import "github.com/oklog/ulid/v2"

// GenerateID returns a new ULID as a string.
func GenerateID() string {
	return ulid.Make().String()
}
