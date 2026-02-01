package internal

import "testing"

func BenchmarkUUIDsStrings1(b *testing.B) {
	b.ReportAllocs()

	uuids := make(UUIDs, 1)
	uuids[0] = New()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = uuids.Strings()
	}
}

func BenchmarkUUIDsStrings10(b *testing.B) {
	b.ReportAllocs()

	uuids := make(UUIDs, 10)
	for i := range uuids {
		uuids[i] = New()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = uuids.Strings()
	}
}

func BenchmarkUUIDsStrings100(b *testing.B) {
	b.ReportAllocs()

	uuids := make(UUIDs, 100)
	for i := range uuids {
		uuids[i] = New()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = uuids.Strings()
	}
}

func BenchmarkUUIDsStrings1000(b *testing.B) {
	b.ReportAllocs()

	uuids := make(UUIDs, 1000)
	for i := range uuids {
		uuids[i] = New()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = uuids.Strings()
	}
}

func BenchmarkUUIDsStringsEmpty(b *testing.B) {
	b.ReportAllocs()

	var uuids UUIDs

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = uuids.Strings()
	}
}
