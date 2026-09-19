package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2Password hashes passwords using Argon2id.
type Argon2Password struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	KeyLength   uint32
	SaltLength  int
}

var _ Password = (*Argon2Password)(nil)

// NewArgon2Password creates an Argon2id password hasher with secure defaults.
func NewArgon2Password() *Argon2Password {
	return &Argon2Password{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 4,
		KeyLength:   32,
		SaltLength:  16,
	}
}

func (p *Argon2Password) Hash(password string) (string, error) {
	salt := make([]byte, p.SaltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("generate password salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		p.Iterations,
		p.Memory,
		p.Parallelism,
		p.KeyLength,
	)

	return fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		p.Memory,
		p.Iterations,
		p.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

func (*Argon2Password) Verify(encodedHash, password string) bool {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" || parts[2] != "v=19" {
		return false
	}

	memory, iterations, parallelism, ok := parseArgon2Parameters(parts[3])
	if !ok {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil || len(salt) == 0 {
		return false
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil || len(expectedHash) == 0 {
		return false
	}

	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		uint32(len(expectedHash)),
	)

	return subtle.ConstantTimeCompare(actualHash, expectedHash) == 1
}

func parseArgon2Parameters(encoded string) (uint32, uint32, uint8, bool) {
	var memory, iterations uint32
	var parallelism uint8

	for _, parameter := range strings.Split(encoded, ",") {
		keyAndValue := strings.SplitN(parameter, "=", 2)
		if len(keyAndValue) != 2 {
			return 0, 0, 0, false
		}

		switch keyAndValue[0] {
		case "m":
			value, err := strconv.ParseUint(keyAndValue[1], 10, 32)
			if err != nil {
				return 0, 0, 0, false
			}
			memory = uint32(value)
		case "t":
			value, err := strconv.ParseUint(keyAndValue[1], 10, 32)
			if err != nil {
				return 0, 0, 0, false
			}
			iterations = uint32(value)
		case "p":
			value, err := strconv.ParseUint(keyAndValue[1], 10, 8)
			if err != nil {
				return 0, 0, 0, false
			}
			parallelism = uint8(value)
		default:
			return 0, 0, 0, false
		}
	}

	if memory == 0 || iterations == 0 || parallelism == 0 {
		return 0, 0, 0, false
	}

	return memory, iterations, parallelism, true
}
