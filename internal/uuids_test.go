package internal

import "testing"

func TestUUIDsStrings(t *testing.T) {
	uuids := UUIDs{
		MustParse("12345678-1234-1234-1234-123456789012"),
		MustParse("abcdefab-cdef-abcd-efab-cdefabcdefab"),
	}

	strings := uuids.Strings()

	if len(strings) != 2 {
		t.Errorf("Strings() len = %d, want 2", len(strings))
	}

	if strings[0] != "12345678-1234-1234-1234-123456789012" {
		t.Errorf("Strings()[0] = %q, want 12345678-1234-1234-1234-123456789012", strings[0])
	}

	if strings[1] != "abcdefab-cdef-abcd-efab-cdefabcdefab" {
		t.Errorf("Strings()[1] = %q, want abcdefab-cdef-abcd-efab-cdefabcdefab", strings[1])
	}
}

func TestUUIDsStringsEmpty(t *testing.T) {
	var uuids UUIDs

	strings := uuids.Strings()

	if strings == nil {
		t.Error("Strings() should return empty slice, not nil")
	}

	if len(strings) != 0 {
		t.Errorf("Strings() len = %d, want 0", len(strings))
	}
}

func TestUUIDsStringsNil(t *testing.T) {
	uuids := UUIDs(nil)
	strings := uuids.Strings()

	if strings == nil {
		t.Error("Strings() should return empty slice, not nil")
	}

	if len(strings) != 0 {
		t.Errorf("Strings() len = %d, want 0", len(strings))
	}
}

func TestUUIDsSliceOperations(t *testing.T) {
	uuids := make(UUIDs, 0, 10)
	for i := 0; i < 5; i++ {
		uuids = append(uuids, New())
	}

	if len(uuids) != 5 {
		t.Errorf("len(uuids) = %d, want 5", len(uuids))
	}

	if cap(uuids) != 10 {
		t.Errorf("cap(uuids) = %d, want 10", cap(uuids))
	}
}

func TestUUIDsIndexing(t *testing.T) {
	uuid1 := MustParse("12345678-1234-1234-1234-123456789012")
	uuid2 := MustParse("abcdefab-cdef-abcd-efab-cdefabcdefab")

	uuids := UUIDs{uuid1, uuid2}

	if uuids[0] != uuid1 {
		t.Errorf("uuids[0] = %v, want %v", uuids[0], uuid1)
	}

	if uuids[1] != uuid2 {
		t.Errorf("uuids[1] = %v, want %v", uuids[1], uuid2)
	}
}
