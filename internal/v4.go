package internal

import (
	"crypto/rand"
	"io"
)

// New generates a new random UUID v4.
func New() UUID {
	var uuid UUID

	err := uuid.readRandom()
	if err != nil {
		panic(err)
	}

	return uuid
}

// NewString generates a new random UUID v4 string using in-place conversion.
func NewString() string {
	return New().String()
}

// NewRandom generates a new random UUID v4, returning an error if crypto/rand fails.
func NewRandom() (UUID, error) {
	var uuid UUID

	err := uuid.readRandom()
	if err != nil {
		return Nil, err
	}

	return uuid, nil
}

// NewRandomFromReader generates a new random UUID v4 from a reader.
func NewRandomFromReader(r io.Reader) (UUID, error) {
	var uuid UUID

	_, err := r.Read(uuid[:])
	if err != nil {
		return Nil, err
	}

	uuid.SetVersion(V4)
	uuid.SetVariant(VariantRFC4122)

	return uuid, nil
}

func (u *UUID) readRandom() error {
	_, err := rand.Read(u[:])
	if err != nil {
		return err
	}

	u.SetVersion(V4)
	u.SetVariant(VariantRFC4122)

	return nil
}
