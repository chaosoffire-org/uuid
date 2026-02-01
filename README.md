# ⚡ uuid

<div align="center">

[![Go Reference](https://pkg.go.dev/badge/github.com/chaosoffire-org/uuid.svg)](https://pkg.go.dev/github.com/chaosoffire-org/uuid)
[![Go Report Card](https://goreportcard.com/badge/github.com/chaosoffire-org/uuid)](https://goreportcard.com/report/github.com/chaosoffire-org/uuid)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Coverage](https://img.shields.io/badge/Coverage-98.4%25-brightgreen.svg)](coverage.html)

**A blazing-fast, zero-allocation UUID library for Go**

_RFC 9562 compliant • Outperforms google/uuid • Battle-tested with fuzz testing_

</div>

---

## ✨ Highlights

| Feature                         | Description                              |
| ------------------------------- | ---------------------------------------- |
| 🚀 **Zero-Allocation Parsing**  | Parse UUIDs without heap allocations     |
| ⚡ **Significantly Faster**     | Optimized for speed vs google/uuid       |
| 🔒 **Cryptographically Secure** | Uses `crypto/rand` for random generation |
| 📦 **Zero Dependencies**        | No external dependencies whatsoever      |
| 🧪 **98.4% Test Coverage**      | Comprehensive unit + fuzz testing        |
| 🛡️ **Bounds Check Elimination** | Compiler-verified safe array access      |

## 🎯 Supported Versions

| Version | Description                      | Status         |
| ------- | -------------------------------- | -------------- |
| **V4**  | Random-based UUID                | ✅ Implemented |
| **V7**  | Unix Epoch time-based (RFC 9562) | ✅ Implemented |

---

## 📦 Installation

```bash
go get github.com/chaosoffire-org/uuid
```

## 🚀 Quick Start

```go
package main

import (
    "fmt"
    "github.com/chaosoffire-org/uuid"
)

func main() {
    // Generate a random UUID v4
    id := uuid.New()
    fmt.Println(id.String()) // "6ba7b810-9dad-4d1a-80b4-00c04fd430c8"

    // Generate a time-ordered UUID v7 (monotonically increasing)
    id7, _ := uuid.NewV7()
    fmt.Println(id7.String()) // "019475fc-9c58-7b34-8b9a-4d9e3f2c1a5b"

    // Parse a UUID string
    parsed, err := uuid.Parse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
    if err != nil {
        panic(err)
    }
    fmt.Println(parsed.Version()) // VERSION_1
}
```

---

## 📚 API Reference

### Generation

```go
// V4 (Random)
uuid.New()              // UUID (panics on crypto/rand failure)
uuid.NewRandom()        // (UUID, error)
uuid.NewString()        // string (direct, no intermediate UUID)

// V7 (Time-based, monotonically increasing)
uuid.NewV7()            // (UUID, error)
```

### Parsing

```go
// All formats supported
uuid.Parse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")           // Standard
uuid.Parse("urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8")  // URN
uuid.Parse("{6ba7b810-9dad-11d1-80b4-00c04fd430c8}")         // Braced
uuid.Parse("6ba7b8109dad11d180b400c04fd430c8")               // No dashes

uuid.ParseBytes([]byte("..."))  // From []byte
uuid.MustParse("...")           // Panics on error
uuid.Validate("...")            // Returns error if invalid
```

### UUID Methods

```go
id.String()          // "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
id.URN()             // "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8"
id.Version()         // uuid.V4
id.Variant()         // uuid.VariantRFC4122
id.Bytes()           // []byte (16 bytes, zero-copy)
id.MarshalText()     // encoding.TextMarshaler
id.MarshalBinary()   // encoding.BinaryMarshaler
```

### Database Integration

Implements `sql.Scanner` and `driver.Valuer`:

```go
type User struct {
    ID   uuid.UUID `db:"id"`
    Name string    `db:"name"`
}

// Works seamlessly with database/sql
db.QueryRow("SELECT id, name FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Name)
```

### Nullable UUID

```go
var nullID uuid.NullUUID
db.QueryRow("SELECT nullable_id FROM users").Scan(&nullID)
if nullID.Valid {
    fmt.Println(nullID.UUID)
}
```

---

## 📊 Benchmarks

Tested on Go 1.25.

Run the benchmarks yourself to verify performance:

```bash
make benchmark-compare
```

This command runs a comprehensive comparison against `google/uuid`, measuring parsing speed, generation time, and memory allocations.

---

## ⚠️ Design Philosophy: `unsafe` Usage

This library uses `unsafe` extensively for **maximum performance**. Every usage is:

- ✅ **Carefully bounded** – All pointer arithmetic uses compile-time known offsets
- ✅ **Fuzz tested** – Validated with millions of random inputs
- ✅ **Bounds-check eliminated** – Compiler verifies safety statically

### Compiler Verification

You can verify the optimizations (inlining and bounds-check elimination) by running:

```bash
go build -gcflags='-m -m' ./internal 2>&1 | grep -E "index bounds check elided|can inline"
```

> [!NOTE]
> All `unsafe` usages follow Go's [unsafe.Pointer rules](https://pkg.go.dev/unsafe#Pointer).
> The codebase is regularly validated with `go vet` and `golangci-lint`.

---

## 🧪 Testing

```bash
make test              # Run unit tests
make test-cover        # Generate coverage report
make benchmark         # Run all benchmarks
make benchmark-compare # Compare vs google/uuid
make fuzz              # Run all fuzz tests
make fuzz-list         # List available fuzz targets
make fuzz-single NAME=FuzzParse  # Run specific fuzz test
```

### Test Coverage

The library maintains high test coverage, verified by `codecov`.

_Note: Uncovered code is primarily `crypto/rand` failure paths (untestable without mocking)_

---

## 🔮 Why Not V1/V3/V5?

| Version   | Reason                                                 |
| --------- | ------------------------------------------------------ |
| **V1**    | Uses MAC address (privacy concern), superseded by V7   |
| **V3/V5** | Hash-based, rarely needed; trivial to add if requested |
| **V7**    | Modern RFC 9562 standard for time-ordered UUIDs ✅     |

---

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING](CONTRIBUTING.md) for our CLA and guidelines.

```bash
make format          # Format code
make lint            # Check for lint errors
make test            # Run unit tests
make benchmark       # Ensure no performance regression
make fuzz            # Run fuzz tests
```

---

## 📄 License

Please see [LICENSE](LICENSE) for details.

---

## 🙏 Acknowledgments

- Inspired by [google/uuid](https://github.com/google/uuid)
- UUID specification: [RFC 9562](https://www.rfc-editor.org/rfc/rfc9562.html)
- Performance techniques from the Go standard library
