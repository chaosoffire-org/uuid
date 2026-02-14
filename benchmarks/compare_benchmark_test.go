package benchmarks

import (
	"testing"

	"github.com/chaosoffire-org/uuid"

	googleuuid "github.com/google/uuid"
)

// =============================================================================
// UUID Generation Benchmarks
// =============================================================================

var (
	sinkUUID          uuid.UUID
	sinkGoogleUUID    googleuuid.UUID
	sinkUUIDs         []string
	sinkGoogleUUIDs   []string
	sinkString        string
	sinkBytes         []byte
	sinkErr           error
	sinkVersion       uuid.Version
	sinkGoogleVersion googleuuid.Version
	sinkVariant       uuid.Variant
	sinkGoogleVariant googleuuid.Variant
)

func BenchmarkNew_This(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkUUID = uuid.New()
	}
}

func BenchmarkNew_Google(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkGoogleUUID = googleuuid.New()
	}
}

func BenchmarkNewRandom_This(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkUUID, sinkErr = uuid.NewRandom()
	}
}

func BenchmarkNewRandom_Google(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkGoogleUUID, sinkErr = googleuuid.NewRandom()
	}
}

func BenchmarkNewString_This(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkString = uuid.NewString()
	}
}

func BenchmarkNewString_Google(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkString = googleuuid.NewString()
	}
}

// =============================================================================
// UUID Parsing Benchmarks
// =============================================================================

func BenchmarkParse_This(b *testing.B) {
	b.ReportAllocs()
	s := "12345678-1234-1234-1234-123456789012"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkUUID, sinkErr = uuid.Parse(s)
	}
}

func BenchmarkParse_Google(b *testing.B) {
	b.ReportAllocs()
	s := "12345678-1234-1234-1234-123456789012"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkGoogleUUID, sinkErr = googleuuid.Parse(s)
	}
}

func BenchmarkParseBytes_This(b *testing.B) {
	b.ReportAllocs()
	data := []byte("12345678-1234-1234-1234-123456789012")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkUUID, sinkErr = uuid.ParseBytes(data)
	}
}

func BenchmarkParseBytes_Google(b *testing.B) {
	b.ReportAllocs()
	data := []byte("12345678-1234-1234-1234-123456789012")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkGoogleUUID, sinkErr = googleuuid.ParseBytes(data)
	}
}

func BenchmarkParseURN_This(b *testing.B) {
	b.ReportAllocs()
	s := "urn:uuid:12345678-1234-1234-1234-123456789012"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkUUID, sinkErr = uuid.Parse(s)
	}
}

func BenchmarkParseURN_Google(b *testing.B) {
	b.ReportAllocs()
	s := "urn:uuid:12345678-1234-1234-1234-123456789012"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkGoogleUUID, sinkErr = googleuuid.Parse(s)
	}
}

func BenchmarkParseBraces_This(b *testing.B) {
	b.ReportAllocs()
	s := "{12345678-1234-1234-1234-123456789012}"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkUUID, sinkErr = uuid.Parse(s)
	}
}

func BenchmarkParseBraces_Google(b *testing.B) {
	b.ReportAllocs()
	s := "{12345678-1234-1234-1234-123456789012}"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkGoogleUUID, sinkErr = googleuuid.Parse(s)
	}
}

func BenchmarkParseNoDashes_This(b *testing.B) {
	b.ReportAllocs()
	s := "12345678123412341234123456789012"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkUUID, sinkErr = uuid.Parse(s)
	}
}

func BenchmarkParseNoDashes_Google(b *testing.B) {
	b.ReportAllocs()
	s := "12345678123412341234123456789012"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkGoogleUUID, sinkErr = googleuuid.Parse(s)
	}
}

// =============================================================================
// UUID String Conversion Benchmarks
// =============================================================================

func BenchmarkString_This(b *testing.B) {
	b.ReportAllocs()
	u := uuid.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkString = u.String()
	}
}

func BenchmarkString_Google(b *testing.B) {
	b.ReportAllocs()
	u := googleuuid.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkString = u.String()
	}
}

func BenchmarkURN_This(b *testing.B) {
	b.ReportAllocs()
	u := uuid.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkString = u.URN()
	}
}

func BenchmarkURN_Google(b *testing.B) {
	b.ReportAllocs()
	u := googleuuid.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkString = u.URN()
	}
}

func BenchmarkStrings_This(b *testing.B) {
	b.ReportAllocs()
	const count = 1000
	uuids := make(uuid.UUIDs, count)
	for i := 0; i < count; i++ {
		uuids[i] = uuid.New()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkUUIDs = uuids.Strings()
	}
}

func BenchmarkStrings_Google(b *testing.B) {
	b.ReportAllocs()
	const count = 1000
	uuids := make(googleuuid.UUIDs, count)
	for i := 0; i < count; i++ {
		uuids[i] = googleuuid.New()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkGoogleUUIDs = uuids.Strings()
	}
}

// =============================================================================
// UUID Marshaling Benchmarks
// =============================================================================

func BenchmarkMarshalText_This(b *testing.B) {
	b.ReportAllocs()
	u := uuid.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkBytes, sinkErr = u.MarshalText()
	}
}

func BenchmarkMarshalText_Google(b *testing.B) {
	b.ReportAllocs()
	u := googleuuid.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkBytes, sinkErr = u.MarshalText()
	}
}

func BenchmarkUnmarshalText_This(b *testing.B) {
	b.ReportAllocs()
	var u uuid.UUID
	data := []byte("12345678-1234-1234-1234-123456789012")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkErr = u.UnmarshalText(data)
	}
}

func BenchmarkUnmarshalText_Google(b *testing.B) {
	b.ReportAllocs()
	var u googleuuid.UUID
	data := []byte("12345678-1234-1234-1234-123456789012")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkErr = u.UnmarshalText(data)
	}
}

func BenchmarkMarshalBinary_This(b *testing.B) {
	b.ReportAllocs()
	u := uuid.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkBytes, sinkErr = u.MarshalBinary()
	}
}

func BenchmarkMarshalBinary_Google(b *testing.B) {
	b.ReportAllocs()
	u := googleuuid.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkBytes, sinkErr = u.MarshalBinary()
	}
}

func BenchmarkUnmarshalBinary_This(b *testing.B) {
	b.ReportAllocs()
	var u uuid.UUID
	data := make([]byte, 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkErr = u.UnmarshalBinary(data)
	}
}

func BenchmarkUnmarshalBinary_Google(b *testing.B) {
	b.ReportAllocs()
	var u googleuuid.UUID
	data := make([]byte, 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkErr = u.UnmarshalBinary(data)
	}
}

// =============================================================================
// UUID Version/Variant Benchmarks
// =============================================================================

func BenchmarkVersion_This(b *testing.B) {
	b.ReportAllocs()
	u := uuid.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkVersion = u.Version()
	}
}

func BenchmarkVersion_Google(b *testing.B) {
	b.ReportAllocs()
	u := googleuuid.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkGoogleVersion = u.Version()
	}
}

func BenchmarkVariant_This(b *testing.B) {
	b.ReportAllocs()
	u := uuid.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkVariant = u.Variant()
	}
}

func BenchmarkVariant_Google(b *testing.B) {
	b.ReportAllocs()
	u := googleuuid.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkGoogleVariant = u.Variant()
	}
}

// =============================================================================
// UUID V7 Generation Benchmarks
// =============================================================================

func BenchmarkNewV7_This(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkUUID, sinkErr = uuid.NewV7()
	}
}

func BenchmarkNewV7_Google(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkGoogleUUID, sinkErr = googleuuid.NewV7()
	}
}

// =============================================================================
// FromBytes Benchmarks
// =============================================================================

func BenchmarkFromBytes_This(b *testing.B) {
	b.ReportAllocs()
	data := make([]byte, 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkUUID, sinkErr = uuid.FromBytes(data)
	}
}

func BenchmarkFromBytes_Google(b *testing.B) {
	b.ReportAllocs()
	data := make([]byte, 16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkGoogleUUID, sinkErr = googleuuid.FromBytes(data)
	}
}
