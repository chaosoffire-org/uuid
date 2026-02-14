package internal

import (
	"testing"
)

func TestEncodeHex(t *testing.T) {
	uuid := UUID{0x12, 0x34, 0x56, 0x78, 0x12, 0x34, 0x12, 0x34, 0x12, 0x34, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12}
	expected := "12345678-1234-1234-1234-123456789012"

	var buf [36]byte
	encodeHex(&buf, uuid)

	result := string(buf[:])
	if result != expected {
		t.Errorf("encodeHex() = %q, want %q", result, expected)
	}
}

func TestEncodeHexDashPositions(t *testing.T) {
	uuid := New()

	var buf [36]byte
	encodeHex(&buf, uuid)

	// Verify dash positions
	if buf[8] != '-' || buf[13] != '-' || buf[18] != '-' || buf[23] != '-' {
		t.Errorf("encodeHex() dashes at wrong positions: %s", string(buf[:]))
	}
}

func TestEncodeHexAllBytes(t *testing.T) {
	// Test with all byte values
	for b := 0; b < 256; b++ {
		uuid := UUID{}
		uuid[0] = byte(b)

		var buf [36]byte
		encodeHex(&buf, uuid)

		// First two characters should represent the byte
		high := hextable[byte(b)>>4]

		low := hextable[byte(b)&0x0f]
		if buf[0] != high || buf[1] != low {
			t.Errorf("encodeHex with byte 0x%02x: got %c%c, want %c%c", b, buf[0], buf[1], high, low)
		}
	}
}

func TestReadRandom(t *testing.T) {
	var dst [16]byte

	err := readRandom(&dst)
	if err != nil {
		t.Errorf("readRandom() error = %v, wantErr false", err)
	}

	// Check that the bytes are not all zero
	for _, b := range dst {
		if b != 0 {
			return
		}
	}

	t.Errorf("readRandom() returned all zeros")
}

func TestReadRandomUniqueness(t *testing.T) {
	var a, b [16]byte

	if err := readRandom(&a); err != nil {
		t.Fatalf("readRandom() error = %v", err)
	}

	if err := readRandom(&b); err != nil {
		t.Fatalf("readRandom() error = %v", err)
	}

	if a == b {
		t.Errorf("readRandom() produced identical outputs: %x", a)
	}
}
