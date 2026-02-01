package internal

import (
	"testing"
)

func TestNullUUIDScanNil(t *testing.T) {
	var nu NullUUID

	err := nu.Scan(nil)
	if err != nil {
		t.Errorf("Scan(nil) returned error: %v", err)
	}

	if nu.Valid {
		t.Error("Scan(nil) should set Valid to false")
	}

	if nu.UUID != Nil {
		t.Error("Scan(nil) should set UUID to Nil")
	}
}

func TestNullUUIDScanString(t *testing.T) {
	var nu NullUUID

	err := nu.Scan("12345678-1234-1234-1234-123456789012")
	if err != nil {
		t.Errorf("Scan(string) returned error: %v", err)
	}

	if !nu.Valid {
		t.Error("Scan(string) should set Valid to true")
	}

	expected := MustParse("12345678-1234-1234-1234-123456789012")
	if nu.UUID != expected {
		t.Errorf("Scan(string) UUID = %v, want %v", nu.UUID, expected)
	}
}

func TestNullUUIDScanBytes(t *testing.T) {
	var nu NullUUID

	input := []byte("12345678-1234-1234-1234-123456789012")

	err := nu.Scan(input)
	if err != nil {
		t.Errorf("Scan([]byte) returned error: %v", err)
	}

	if !nu.Valid {
		t.Error("Scan([]byte) should set Valid to true")
	}
}

func TestNullUUIDScanInvalid(t *testing.T) {
	var nu NullUUID

	err := nu.Scan("invalid")
	if err == nil {
		t.Error("Scan(invalid) should return error")
	}
}

func TestNullUUIDScanUnsupportedType(t *testing.T) {
	var nu NullUUID

	err := nu.Scan(12345)
	if err != ErrUnsupportedScanType {
		t.Errorf("Scan(int) error = %v, want ErrUnsupportedScanType", err)
	}
}

func TestNullUUIDValueValid(t *testing.T) {
	uuid := MustParse("12345678-1234-1234-1234-123456789012")
	nu := NullUUID{UUID: uuid, Valid: true}

	val, err := nu.Value()
	if err != nil {
		t.Errorf("Value() returned error: %v", err)
	}

	if val == nil {
		t.Error("Value() should not be nil for valid NullUUID")
	}

	str, ok := val.(string)
	if !ok {
		t.Errorf("Value() should return a string type")
	}

	if str != "12345678-1234-1234-1234-123456789012" {
		t.Errorf("Value() = %v, want 12345678-1234-1234-1234-123456789012", str)
	}
}

func TestNullUUIDValueInvalid(t *testing.T) {
	nu := NullUUID{Valid: false}

	val, err := nu.Value()
	if err != nil {
		t.Errorf("Value() returned error: %v", err)
	}

	if val != nil {
		t.Errorf("Value() for invalid NullUUID = %v, want nil", val)
	}
}

func TestNullUUIDZeroValue(t *testing.T) {
	var nu NullUUID
	if nu.Valid {
		t.Error("zero value NullUUID should have Valid = false")
	}

	if nu.UUID != Nil {
		t.Error("zero value NullUUID should have UUID = Nil")
	}
}
