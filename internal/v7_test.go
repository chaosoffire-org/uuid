package internal

import (
	"testing"
	"time"
)

func TestNewV7(t *testing.T) {
	uuid, err := NewV7()
	if err != nil {
		t.Fatalf("NewV7() error = %v", err)
	}

	if uuid.Version() != V7 {
		t.Errorf("NewV7() Version = %v, want V7", uuid.Version())
	}

	if uuid.Variant() != VariantRFC4122 {
		t.Errorf("NewV7() Variant = %v, want VariantRFC4122", uuid.Variant())
	}
}

func TestNewV7_Timestamp(t *testing.T) {
	before := time.Now().UnixMilli()
	uuid, err := NewV7()
	after := time.Now().UnixMilli()

	if err != nil {
		t.Fatalf("NewV7() error = %v", err)
	}

	// Extract timestamp from UUID (first 48 bits)
	ts := uint64(uuid[0])<<40 | uint64(uuid[1])<<32 | uint64(uuid[2])<<24 |
		uint64(uuid[3])<<16 | uint64(uuid[4])<<8 | uint64(uuid[5])

	if ts < uint64(before) || ts > uint64(after)+1000 {
		t.Errorf("NewV7() timestamp = %d, want between %d and %d", ts, before, after+1000)
	}
}

func TestNewV7_Monotonic(t *testing.T) {
	const count = 10000

	uuids := make([]UUID, count)

	for i := 0; i < count; i++ {
		uuid, err := NewV7()
		if err != nil {
			t.Fatalf("NewV7() error = %v", err)
		}

		uuids[i] = uuid
	}

	// Check strict monotonicity
	for i := 1; i < count; i++ {
		if compareV7(uuids[i-1], uuids[i]) >= 0 {
			t.Errorf("UUIDs not strictly increasing at index %d: %s >= %s",
				i, uuids[i-1].String(), uuids[i].String())
		}
	}
}

// compareV7 compares two V7 UUIDs for ordering.
// Returns -1 if a < b, 0 if a == b, 1 if a > b.
func compareV7(a, b UUID) int {
	for i := 0; i < 16; i++ {
		if a[i] < b[i] {
			return -1
		}

		if a[i] > b[i] {
			return 1
		}
	}

	return 0
}

func BenchmarkNewV7(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = NewV7()
	}
}
