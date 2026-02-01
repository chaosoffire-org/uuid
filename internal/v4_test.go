package internal

import (
	"errors"
	"io"
	"testing"
)

// errorReader is a mock reader that always returns an error.
type errorReader struct{}

func (r errorReader) Read(p []byte) (int, error) {
	return 0, errors.New("mock read error")
}

// limitedReader reads only a limited number of bytes before returning EOF.
type limitedReader struct {
	remaining int
}

func (r *limitedReader) Read(p []byte) (int, error) {
	if r.remaining <= 0 {
		return 0, io.EOF
	}

	n := len(p)
	if n > r.remaining {
		n = r.remaining
	}

	r.remaining -= n

	return n, nil
}

func TestNew(t *testing.T) {
	uuid := New()
	if uuid == Nil {
		t.Error("New() returned Nil UUID")
	}

	if uuid.Version() != V4 {
		t.Errorf("New() Version = %v, want V4", uuid.Version())
	}

	if uuid.Variant() != VariantRFC4122 {
		t.Errorf("New() Variant = %v, want VariantRFC4122", uuid.Variant())
	}
}

func TestNewUniqueness(t *testing.T) {
	seen := make(map[UUID]bool)

	for i := 0; i < 1000; i++ {
		uuid := New()
		if seen[uuid] {
			t.Error("New() returned duplicate UUID")
		}

		seen[uuid] = true
	}
}

func TestNewString(t *testing.T) {
	s := NewString()
	if len(s) != 36 {
		t.Errorf("NewString() len = %d, want 36", len(s))
	}
	// Should be parseable
	_, err := Parse(s)
	if err != nil {
		t.Errorf("NewString() not parseable: %v", err)
	}
}

func TestNewRandom(t *testing.T) {
	uuid, err := NewRandom()
	if err != nil {
		t.Errorf("NewRandom() error = %v", err)
	}

	if uuid == Nil {
		t.Error("NewRandom() returned Nil UUID")
	}

	if uuid.Version() != V4 {
		t.Errorf("NewRandom() Version = %v, want V4", uuid.Version())
	}
}

func TestNewRandomFromReader(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Create a reader with sufficient bytes
		reader := &limitedReader{remaining: 16}

		uuid, err := NewRandomFromReader(reader)
		if err != nil {
			t.Errorf("NewRandomFromReader() error = %v", err)
		}

		if uuid.Version() != V4 {
			t.Errorf("NewRandomFromReader() Version = %v, want V4", uuid.Version())
		}

		if uuid.Variant() != VariantRFC4122 {
			t.Errorf("NewRandomFromReader() Variant = %v, want VariantRFC4122", uuid.Variant())
		}
	})

	t.Run("error", func(t *testing.T) {
		// Use errorReader to trigger the error path
		uuid, err := NewRandomFromReader(errorReader{})
		if err == nil {
			t.Error("NewRandomFromReader() expected error, got nil")
		}

		if uuid != Nil {
			t.Errorf("NewRandomFromReader() with error should return Nil, got %v", uuid)
		}
	})
}

func TestNewRandomUniqueness(t *testing.T) {
	seen := make(map[UUID]bool)

	for i := 0; i < 1000; i++ {
		uuid, err := NewRandom()
		if err != nil {
			t.Fatalf("NewRandom() error = %v", err)
		}

		if seen[uuid] {
			t.Error("NewRandom() returned duplicate UUID")
		}

		seen[uuid] = true
	}
}
