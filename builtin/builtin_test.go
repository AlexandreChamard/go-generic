package builtin

import "testing"

func TestLess(t *testing.T) {
	Less(uint8(0), uint8(0))
	Less(uint16(0), uint16(0))
	Less(uint32(0), uint32(0))
	Less(uint64(0), uint64(0))
	Less(int8(0), int8(0))
	Less(int16(0), int16(0))
	Less(int32(0), int32(0))
	Less(int64(0), int64(0))
	Less(float32(0), float32(0))
	Less(float64(0), float64(0))
	Less(string(""), string(""))
	Less(int(0), int(0))
	Less(uint(0), uint(0))
	Less(uintptr(0), uintptr(0))
}
