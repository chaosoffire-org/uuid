package internal

import "testing"

func BenchmarkHex(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for n := byte(0); n < 16; n++ {
			_ = hextable[n]
		}
	}
}

func BenchmarkHexSingle(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = hextable[10]
	}
}

func BenchmarkEncodeHex(b *testing.B) {
	b.ReportAllocs()

	uuid := New()

	var buf [36]byte

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		encodeHex(&buf, uuid)
	}
}
