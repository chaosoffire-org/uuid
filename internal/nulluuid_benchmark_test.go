package internal

import "testing"

func BenchmarkNullUUIDScanNil(b *testing.B) {
	b.ReportAllocs()

	var nu NullUUID
	for i := 0; i < b.N; i++ {
		_ = nu.Scan(nil)
	}
}

func BenchmarkNullUUIDScanString(b *testing.B) {
	b.ReportAllocs()

	var nu NullUUID

	s := "12345678-1234-1234-1234-123456789012"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = nu.Scan(s)
	}
}

func BenchmarkNullUUIDScanBytes(b *testing.B) {
	b.ReportAllocs()

	var nu NullUUID

	data := []byte("12345678-1234-1234-1234-123456789012")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = nu.Scan(data)
	}
}

func BenchmarkNullUUIDScanRawBytes(b *testing.B) {
	b.ReportAllocs()

	var nu NullUUID

	uuid := New()
	data := uuid[:]

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = nu.Scan(data)
	}
}

func BenchmarkNullUUIDValueValid(b *testing.B) {
	b.ReportAllocs()

	nu := NullUUID{UUID: New(), Valid: true}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = nu.Value()
	}
}

func BenchmarkNullUUIDValueInvalid(b *testing.B) {
	b.ReportAllocs()

	nu := NullUUID{Valid: false}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = nu.Value()
	}
}
