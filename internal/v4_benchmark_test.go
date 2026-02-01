package internal

import "testing"

func BenchmarkNew(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = New()
	}
}

func BenchmarkNewString(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = NewString()
	}
}

func BenchmarkNewRandom(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = NewRandom()
	}
}
