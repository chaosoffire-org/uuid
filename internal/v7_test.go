package internal

import (
	"sync"
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

func TestNewV7_Concurrent(t *testing.T) {
	const (
		goroutines = 100
		perRoutine = 100
	)

	ch := make(chan UUID, goroutines*perRoutine)

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()

			for i := 0; i < perRoutine; i++ {
				uuid, err := NewV7()
				if err != nil {
					t.Errorf("NewV7() error = %v", err)
					return
				}

				ch <- uuid
			}
		}()
	}

	wg.Wait()
	close(ch)

	if len(ch) != goroutines*perRoutine {
		t.Errorf("expected %d UUIDs, got %d", goroutines*perRoutine, len(ch))
	}

	seen := make(map[UUID]struct{})
	for uuid := range ch {
		seen[uuid] = struct{}{}

		if uuid.Version() != V7 {
			t.Errorf("Version = %v, want V7", uuid.Version())
		}

		if uuid.Variant() != VariantRFC4122 {
			t.Errorf("Variant = %v, want VariantRFC4122", uuid.Variant())
		}
	}

	if len(seen) != goroutines*perRoutine {
		t.Errorf("expected %d unique UUIDs, got %d", goroutines*perRoutine, len(seen))
	}
}

func TestNewV7_SequenceOverflow(t *testing.T) {
	// Lock and set state to force sequence overflow
	v7Mu.Lock()
	origLastTime := v7LastTime
	origSeq := v7Seq

	v7LastTime = uint64(time.Now().UnixMilli()) + 100_000 // far future to avoid real-time collision
	v7Seq = 0x0FFE
	savedTime := v7LastTime
	v7Mu.Unlock()

	t.Cleanup(func() {
		v7Mu.Lock()
		v7LastTime = origLastTime
		v7Seq = origSeq
		v7Mu.Unlock()
	})

	uuid1, err := NewV7()
	if err != nil {
		t.Fatalf("NewV7() error = %v", err)
	}

	uuid2, err := NewV7()
	if err != nil {
		t.Fatalf("NewV7() error = %v", err)
	}

	// uuid2 should trigger overflow: seq goes from 0x0FFF to 0, lastTime increments
	uuid3, err := NewV7()
	if err != nil {
		t.Fatalf("NewV7() error = %v", err)
	}

	// All three should be strictly monotonic
	if compareV7(uuid1, uuid2) >= 0 {
		t.Errorf("uuid1 >= uuid2: %s >= %s", uuid1.String(), uuid2.String())
	}

	if compareV7(uuid2, uuid3) >= 0 {
		t.Errorf("uuid2 >= uuid3: %s >= %s", uuid2.String(), uuid3.String())
	}

	// After overflow, v7LastTime should have incremented
	v7Mu.Lock()
	currentTime := v7LastTime
	v7Mu.Unlock()

	if currentTime <= savedTime {
		t.Errorf("v7LastTime should have incremented after overflow: got %d, started at %d", currentTime, savedTime)
	}
}

func TestNewV7_Uniqueness(t *testing.T) {
	const count = 10000

	seen := make(map[UUID]struct{}, count)
	for i := 0; i < count; i++ {
		uuid, err := NewV7()
		if err != nil {
			t.Fatalf("NewV7() error = %v", err)
		}

		if _, ok := seen[uuid]; ok {
			t.Errorf("duplicate UUID at index %d: %s", i, uuid.String())
		}

		seen[uuid] = struct{}{}
	}
}

func BenchmarkNewV7(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = NewV7()
	}
}
