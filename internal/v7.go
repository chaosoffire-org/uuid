package internal

import (
	"math/bits"
	"sync"
	"time"
	"unsafe"
)

var (
	v7Mu       sync.Mutex
	v7LastTime uint64
	v7Seq      uint16

	isBigEndianVal uint16 = 0x1
	isBigEndian    bool   = *(*byte)(unsafe.Pointer(&isBigEndianVal)) == 0
)

// NewV7 generates a new UUID version 7 (Unix Epoch time-based).
func NewV7() (UUID, error) {
	uuid, err := NewRandom()
	if err != nil {
		return uuid, err
	}

	v7Mu.Lock()

	now := uint64(time.Now().UnixMilli())
	if now > v7LastTime {
		v7LastTime = now
		v7Seq = uint16(uuid[0])<<8 | uint16(uuid[1])&0x0FFF
	} else {
		v7Seq++
		if v7Seq > 0x0FFF {
			v7LastTime++
			v7Seq = 0
		}
	}

	ts := v7LastTime
	seq := v7Seq

	v7Mu.Unlock()

	packed := (ts << 16) | uint64(seq)
	if isBigEndian {
		*(*uint64)(unsafe.Pointer(&uuid[0])) = packed
	} else {
		*(*uint64)(unsafe.Pointer(&uuid[0])) = bits.ReverseBytes64(packed)
	}

	uuid.SetVersion(V7)
	uuid.SetVariant(VariantRFC4122)

	return uuid, nil
}
