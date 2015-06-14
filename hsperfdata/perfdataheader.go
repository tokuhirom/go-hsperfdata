package hsperfdata

import (
	"encoding/binary"
)

type PerfdataHeader struct {
	// 0xcafec0c0
	Magic uint32
	// big_endian == 0, little_endian == 1
	ByteOrder byte
	Major     byte
	Minor     byte
	// Reserved  byte
}

func (header *PerfdataHeader) GetEndian() binary.ByteOrder {
	if header.ByteOrder == 0 {
		return binary.BigEndian
	} else {
		return binary.LittleEndian
	}
}
