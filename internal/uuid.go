// Package internal contains the core UUID implementation.
package internal

import (
	"database/sql/driver"
	"fmt"
	"unsafe"
)

// String returns the string form of uuid.
func (uuid UUID) String() string {
	var buf [36]byte
	encodeHex(&buf, uuid)

	return string(buf[:])
}

// URN returns the RFC 2141 URN form of uuid.
func (uuid UUID) URN() string {
	buf := make([]byte, 36+9)

	copy(buf, "urn:uuid:")
	encodeHex((*[36]byte)(unsafe.Pointer(&buf[9])), uuid)

	return unsafe.String(unsafe.SliceData(buf), 36+9)
}

// Version returns the version of uuid.
func (uuid UUID) Version() Version {
	return Version(uuid[6] >> 4)
}

// Variant returns the variant of uuid.
func (uuid UUID) Variant() Variant {
	switch {
	case (uuid[8] >> 7) == 0x00:
		return VariantNCS
	case (uuid[8] >> 6) == 0x02:
		return VariantRFC4122
	case (uuid[8] >> 5) == 0x06:
		return VariantMicrosoft
	}

	return VariantFuture
}

// SetVersion sets the version bits.
func (uuid *UUID) SetVersion(v Version) {
	uuid[6] = (uuid[6] & 0x0f) | (byte(v) << 4)
}

// SetVariant sets the variant bits.
func (uuid *UUID) SetVariant(v Variant) {
	switch v {
	case VariantNCS:
		uuid[8] = uuid[8] & 0x7f
	case VariantRFC4122:
		uuid[8] = (uuid[8] & 0x3f) | 0x80
	case VariantMicrosoft:
		uuid[8] = (uuid[8] & 0x1f) | 0xc0
	case VariantFuture:
		uuid[8] = (uuid[8] & 0x1f) | 0xe0
	}
}

// Bytes returns the raw bytes.
func (uuid *UUID) Bytes() []byte {
	return uuid[:]
}

// MarshalText implements encoding.TextMarshaler.
func (uuid UUID) MarshalText() ([]byte, error) {
	var buf [36]byte
	encodeHex(&buf, uuid)

	return buf[:], nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (uuid *UUID) UnmarshalText(data []byte) error {
	id, err := ParseBytes(data)
	if err != nil {
		return err
	}

	*uuid = id

	return nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (uuid UUID) MarshalBinary() ([]byte, error) {
	return uuid[:], nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (uuid *UUID) UnmarshalBinary(data []byte) error {
	if len(data) != 16 {
		return invalidLengthError{len(data)}
	}

	*uuid = *(*UUID)(unsafe.Pointer(&data[0]))

	return nil
}

// Value implements driver.Valuer.
func (uuid UUID) Value() (driver.Value, error) {
	return uuid.String(), nil
}

// Scan implements sql.Scanner.
func (uuid *UUID) Scan(src interface{}) error {
	switch src := src.(type) {
	case nil:
		return nil
	case string:
		u, err := Parse(src)
		if err != nil {
			return err
		}

		*uuid = u
	case []byte:
		switch len(src) {
		case 36:
			u, err := ParseBytes(src)
			if err != nil {
				return err
			}

			*uuid = u
		case 16:
			*uuid = *(*UUID)(unsafe.Pointer(&src[0]))
		default:
			return ErrInvalidLength
		}
	default:
		return ErrUnsupportedScanType
	}

	return nil
}

// Parse decodes s into a UUID.
func Parse(s string) (UUID, error) {
	var uuid UUID

	val := unsafe.StringData(s)
	valPtr := unsafe.Pointer(val)

	var ptr unsafe.Pointer

	switch len(s) {
	case 32:
		var acc byte

		ptr = valPtr // 32 bytes
		for i := uintptr(0); i < 16; i++ {
			x1 := *(*byte)(unsafe.Pointer(uintptr(ptr) + i*2))
			x2 := *(*byte)(unsafe.Pointer(uintptr(ptr) + i*2 + 1))

			v1 := xvalues[x1]
			v2 := xvalues[x2]

			acc |= v1 | v2
			uuid[i] = (v1 << 4) | v2
		}

		if acc > 15 {
			return Nil, ErrInvalidUUIDFormat
		}

		return uuid, nil
	case 36:
		ptr = valPtr // 36 bytes
	case 38:
		d1 := *(*byte)(valPtr) ^ '{'
		d2 := *(*byte)(unsafe.Pointer(uintptr(valPtr) + 37)) ^ '}'

		if d1|d2 != 0 {
			return Nil, ErrInvalidBracketedFormat
		}

		ptr = unsafe.Pointer(uintptr(valPtr) + 1) // 36 + 1 bytes
	case 45:
		prefix := unsafe.String((*byte)(valPtr), 9)
		if prefix != "urn:uuid:" {
			return Nil, URNPrefixError{prefix}
		}

		ptr = unsafe.Pointer(uintptr(valPtr) + 9) // 36 bytes
	default:
		return Nil, invalidLengthError{len(s)}
	}

	d1 := *(*byte)(unsafe.Pointer(uintptr(ptr) + 8)) ^ '-'
	d2 := *(*byte)(unsafe.Pointer(uintptr(ptr) + 13)) ^ '-'
	d3 := *(*byte)(unsafe.Pointer(uintptr(ptr) + 18)) ^ '-'
	d4 := *(*byte)(unsafe.Pointer(uintptr(ptr) + 23)) ^ '-'

	if d1|d2|d3|d4 != 0 {
		return Nil, ErrInvalidUUIDFormat
	}

	var acc byte

	// Byte 0: s[0], s[1]
	x1 := *(*byte)(unsafe.Pointer(uintptr(ptr) + 0))
	x2 := *(*byte)(unsafe.Pointer(uintptr(ptr) + 1))

	v1 := xvalues[x1]
	v2 := xvalues[x2]
	acc |= v1 | v2
	uuid[0] = (v1 << 4) | v2

	// Byte 1: s[2], s[3]
	x1 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 2))
	x2 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 3))

	v1 = xvalues[x1]
	v2 = xvalues[x2]
	acc |= v1 | v2
	uuid[1] = (v1 << 4) | v2

	// Byte 2: s[4], s[5]
	x1 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 4))
	x2 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 5))

	v1 = xvalues[x1]
	v2 = xvalues[x2]
	acc |= v1 | v2
	uuid[2] = (v1 << 4) | v2

	// Byte 3: s[6], s[7]
	x1 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 6))
	x2 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 7))

	v1 = xvalues[x1]
	v2 = xvalues[x2]
	acc |= v1 | v2
	uuid[3] = (v1 << 4) | v2

	// s[8] is '-'

	// Byte 4: s[9], s[10]
	x1 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 9))
	x2 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 10))

	v1 = xvalues[x1]
	v2 = xvalues[x2]
	acc |= v1 | v2
	uuid[4] = (v1 << 4) | v2

	// Byte 5: s[11], s[12]
	x1 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 11))
	x2 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 12))

	v1 = xvalues[x1]
	v2 = xvalues[x2]
	acc |= v1 | v2
	uuid[5] = (v1 << 4) | v2

	// s[13] is '-'

	// Byte 6: s[14], s[15]
	x1 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 14))
	x2 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 15))

	v1 = xvalues[x1]
	v2 = xvalues[x2]
	acc |= v1 | v2
	uuid[6] = (v1 << 4) | v2

	// Byte 7: s[16], s[17]
	x1 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 16))
	x2 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 17))

	v1 = xvalues[x1]
	v2 = xvalues[x2]
	acc |= v1 | v2
	uuid[7] = (v1 << 4) | v2

	// s[18] is '-'

	// Byte 8: s[19], s[20]
	x1 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 19))
	x2 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 20))

	v1 = xvalues[x1]
	v2 = xvalues[x2]
	acc |= v1 | v2
	uuid[8] = (v1 << 4) | v2

	// Byte 9: s[21], s[22]
	x1 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 21))
	x2 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 22))

	v1 = xvalues[x1]
	v2 = xvalues[x2]
	acc |= v1 | v2
	uuid[9] = (v1 << 4) | v2

	// s[23] is '-'

	// Byte 10: s[24], s[25]
	x1 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 24))
	x2 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 25))

	v1 = xvalues[x1]
	v2 = xvalues[x2]
	acc |= v1 | v2
	uuid[10] = (v1 << 4) | v2

	// Byte 11: s[26], s[27]
	x1 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 26))
	x2 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 27))

	v1 = xvalues[x1]
	v2 = xvalues[x2]
	acc |= v1 | v2
	uuid[11] = (v1 << 4) | v2

	// Byte 12: s[28], s[29]
	x1 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 28))
	x2 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 29))

	v1 = xvalues[x1]
	v2 = xvalues[x2]
	acc |= v1 | v2
	uuid[12] = (v1 << 4) | v2

	// Byte 13: s[30], s[31]
	x1 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 30))
	x2 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 31))

	v1 = xvalues[x1]
	v2 = xvalues[x2]
	acc |= v1 | v2
	uuid[13] = (v1 << 4) | v2

	// Byte 14: s[32], s[33]
	x1 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 32))
	x2 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 33))

	v1 = xvalues[x1]
	v2 = xvalues[x2]
	acc |= v1 | v2
	uuid[14] = (v1 << 4) | v2

	// Byte 15: s[34], s[35]
	x1 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 34))
	x2 = *(*byte)(unsafe.Pointer(uintptr(ptr) + 35))

	v1 = xvalues[x1]
	v2 = xvalues[x2]
	acc |= v1 | v2
	uuid[15] = (v1 << 4) | v2

	if acc > 15 {
		return Nil, ErrInvalidUUIDFormat
	}

	return uuid, nil
}

// ParseBytes parses a byte slice into a UUID.
func ParseBytes(b []byte) (UUID, error) {
	return Parse(unsafe.String(unsafe.SliceData(b), len(b)))
}

// MustParse is like Parse but panics on error.
func MustParse(s string) UUID {
	uuid, err := Parse(s)
	if err != nil {
		panic(err)
	}

	return uuid
}

// FromBytes creates a UUID from a byte slice.
func FromBytes(b []byte) (UUID, error) {
	if len(b) != 16 {
		return Nil, ErrInvalidLength
	}

	var uuid UUID
	copy(uuid[:], b)

	return uuid, nil
}

// Must panics if err is not nil.
func Must(uuid UUID, err error) UUID {
	if err != nil {
		panic(err)
	}

	return uuid
}

// Version.String returns version as string.
func (v Version) String() string {
	switch v {
	case V1:
		return "VERSION_1"
	case V2:
		return "VERSION_2"
	case V3:
		return "VERSION_3"
	case V4:
		return "VERSION_4"
	case V5:
		return "VERSION_5"
	case V6:
		return "VERSION_6"
	case V7:
		return "VERSION_7"
	}

	return fmt.Sprintf("VERSION_%d", v)
}

// Variant.String returns variant as string.
func (v Variant) String() string {
	switch v {
	case VariantNCS:
		return "NCS"
	case VariantRFC4122:
		return "RFC4122"
	case VariantMicrosoft:
		return "Microsoft"
	}

	return "Future"
}
