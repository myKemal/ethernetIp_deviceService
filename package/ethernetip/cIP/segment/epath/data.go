package epath

import (
	"bytes"
	"strconv"

	"device-ethernetip-go/package/ethernetip/cIP/segment"
	_type "device-ethernetip-go/package/ethernetip/type"
	"device-ethernetip-go/package/lib"
)

type DataTypes _type.USINT

const (
	DataTypeSimple DataTypes = 0x0
	DataTypeANSI   DataTypes = 0x11
	_SIMPLE        DataTypes = 0x80
	_ANSI_EXTD     DataTypes = 0x91
	_UINT8         DataTypes = 0x28
	_UINT16        DataTypes = 0x29
	_UINT32        DataTypes = 0x2a
)

func DataBuild(tp DataTypes, data []byte, padded bool) []byte {
	buffer := new(bytes.Buffer)

	firstByte := uint8(segment.SegmentTypeData) | uint8(tp)
	lib.WriteByte(buffer, firstByte)
	length := uint8(len(data))
	lib.WriteByte(buffer, length)
	lib.WriteByte(buffer, data)
	if padded && buffer.Len()%2 == 1 {
		lib.WriteByte(buffer, uint8(0))
	}

	return buffer.Bytes()
}

func BuildRev(data string) []byte {

	dt, err := strconv.Atoi(data)
	if err != nil {
		return SymbolicBuild(data)
	}
	return ElementBuild(dt)
}

func SymbolicBuild(data string) []byte {
	buffer := new(bytes.Buffer)

	lib.WriteByte(buffer, _ANSI_EXTD)

	lib.WriteByte(buffer, uint8(len(data)))

	lib.WriteByte(buffer, []byte(data))

	if buffer.Len()%2 == 1 {
		lib.WriteByte(buffer, uint8(0))
	}
	return buffer.Bytes()
}

func ElementBuild(data int) []byte {
	buffer := new(bytes.Buffer)

	if data < 256 {
		lib.WriteByte(buffer, _UINT8)
		lib.WriteByte(buffer, uint8(data))
	} else if data < 65536 {
		lib.WriteByte(buffer, _UINT16)
		lib.WriteByte(buffer, uint16(data))
	} else {
		lib.WriteByte(buffer, _UINT32)
		lib.WriteByte(buffer, uint32(data))
	}

	return buffer.Bytes()
}
