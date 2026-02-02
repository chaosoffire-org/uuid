package internal

import "unsafe"

// Strings returns a string slice of all UUIDs.
func (uuids UUIDs) Strings() []string {
	length := len(uuids)
	if length == 0 {
		return []string{}
	}

	buf := make([]byte, length*36)
	result := make([]string, length)

	bufPtr := unsafe.Pointer(unsafe.SliceData(buf))

	for i := 0; i < length; i++ {
		offset := i * 36
		ptr := unsafe.Pointer(uintptr(bufPtr) + uintptr(offset))
		encodeHex((*[36]byte)(ptr), uuids[i])
		result[i] = unsafe.String((*byte)(ptr), 36)
	}

	return result
}
