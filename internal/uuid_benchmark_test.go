package internal

import "testing"

func BenchmarkUUIDString(b *testing.B) {
	b.ReportAllocs()

	uuid := New()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = uuid.String()
	}
}

func BenchmarkUUIDURN(b *testing.B) {
	b.ReportAllocs()

	uuid := New()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = uuid.URN()
	}
}

func BenchmarkUUIDVersion(b *testing.B) {
	b.ReportAllocs()

	uuid := New()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = uuid.Version()
	}
}

func BenchmarkUUIDVariant(b *testing.B) {
	b.ReportAllocs()

	uuid := New()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = uuid.Variant()
	}
}

func BenchmarkSetVersion(b *testing.B) {
	b.ReportAllocs()

	uuid := New()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		uuid.SetVersion(V4)
	}
}

func BenchmarkSetVariant(b *testing.B) {
	b.ReportAllocs()

	uuid := New()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		uuid.SetVariant(VariantRFC4122)
	}
}

func BenchmarkBytes(b *testing.B) {
	b.ReportAllocs()

	uuid := New()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = uuid.Bytes()
	}
}

func BenchmarkMarshalText(b *testing.B) {
	b.ReportAllocs()

	uuid := New()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = uuid.MarshalText()
	}
}

func BenchmarkUnmarshalText(b *testing.B) {
	b.ReportAllocs()

	var uuid UUID

	data := []byte("12345678-1234-1234-1234-123456789012")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = uuid.UnmarshalText(data)
	}
}

func BenchmarkMarshalBinary(b *testing.B) {
	b.ReportAllocs()

	uuid := New()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = uuid.MarshalBinary()
	}
}

func BenchmarkUnmarshalBinary(b *testing.B) {
	b.ReportAllocs()

	var uuid UUID

	data := make([]byte, 16)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = uuid.UnmarshalBinary(data)
	}
}

func BenchmarkValue(b *testing.B) {
	b.ReportAllocs()

	uuid := New()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = uuid.Value()
	}
}

func BenchmarkScanString(b *testing.B) {
	b.ReportAllocs()

	var uuid UUID

	s := "12345678-1234-1234-1234-123456789012"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = uuid.Scan(s)
	}
}

func BenchmarkScanBytes(b *testing.B) {
	b.ReportAllocs()

	var uuid UUID

	data := []byte("12345678-1234-1234-1234-123456789012")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = uuid.Scan(data)
	}
}

func BenchmarkParse(b *testing.B) {
	b.ReportAllocs()

	s := "12345678-1234-1234-1234-123456789012"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Parse(s)
	}
}

func BenchmarkParseURN(b *testing.B) {
	b.ReportAllocs()

	s := "urn:uuid:12345678-1234-1234-1234-123456789012"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Parse(s)
	}
}

func BenchmarkParseBraces(b *testing.B) {
	b.ReportAllocs()

	s := "{12345678-1234-1234-1234-123456789012}"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Parse(s)
	}
}

func BenchmarkParseNoDashes(b *testing.B) {
	b.ReportAllocs()

	s := "12345678123412341234123456789012"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Parse(s)
	}
}

func BenchmarkParseBytes(b *testing.B) {
	b.ReportAllocs()

	data := []byte("12345678-1234-1234-1234-123456789012")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ParseBytes(data)
	}
}

func BenchmarkFromBytes(b *testing.B) {
	b.ReportAllocs()

	data := make([]byte, 16)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = FromBytes(data)
	}
}

func BenchmarkVersionString(b *testing.B) {
	b.ReportAllocs()

	v := V4

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = v.String()
	}
}

func BenchmarkVariantString(b *testing.B) {
	b.ReportAllocs()

	v := VariantRFC4122

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = v.String()
	}
}
