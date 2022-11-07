package tag

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"device-ethernetip-go/package/ethernetip/cIP"
	"device-ethernetip-go/package/ethernetip/cIP/segment/epath"
	"device-ethernetip-go/package/lib"
)

const ServiceReadTag = cIP.Service(0x4c)
const ServiceWriteTag = cIP.Service(0x4d)

type Tag struct {
	name       []byte
	path       []uint8
	readCount  uint16
	xtype      DataType
	structType DataType
	value      []byte
	OnChange   func(interface{})
	OnData     func(interface{})
	next       func()
}

func (t *Tag) GenerateReadMessageRequesRev() *cIP.MessageRouterRequest {
	mr := &cIP.MessageRouterRequest{}
	mr.Service = ServiceReadTag

	pathByte, err := t.GeneratePath()

	if err != nil {
		fmt.Println(err)
		return mr
	}
	mr.RequestPath = pathByte

	data := new(bytes.Buffer)
	lib.WriteByte(data, t.readCount)
	mr.RequestData = data.Bytes()

	return mr
}

func (t *Tag) GenerateReadMessageRequest() *cIP.MessageRouterRequest {
	mr := &cIP.MessageRouterRequest{}
	mr.Service = ServiceReadTag
	mr.RequestPath = epath.DataBuild(epath.DataTypeANSI, t.name, true)

	data := new(bytes.Buffer)
	lib.WriteByte(data, t.readCount)
	mr.RequestData = data.Bytes()

	return mr
}

func (t *Tag) GenerateWriteMessageRequest() *cIP.MessageRouterRequest {

	mr := &cIP.MessageRouterRequest{}
	mr.Service = ServiceWriteTag
	mr.RequestPath = epath.DataBuild(epath.DataTypeANSI, t.name, true)
	data := new(bytes.Buffer)
	lib.WriteByte(data, t.xtype)
	lib.WriteByte(data, t.readCount)
	lib.WriteByte(data, t.GetValue())
	mr.RequestData = data.Bytes()
	return mr
}

func (t *Tag) ReadTagParser(mr *cIP.MessageRouterResponse) {
	fmt.Println("Reading......")

	dataReader := bytes.NewReader(mr.ResponseData)
	lib.ReadByte(dataReader, &t.xtype)
	newValue := make([]byte, dataReader.Len())
	lib.ReadByte(dataReader, newValue)

	if bytes.Compare(t.value, newValue) != 0 {
		t.value = newValue
		if t.OnChange != nil {
			t.OnChange(t.GetValue())
		}
	}

	if t.OnData != nil {
		t.OnData(t.GetValue())
	}

	if t.next != nil {
		f := t.next
		t.next = nil
		f()
	}
}

func (t *Tag) WriteTagParser(mr *cIP.MessageRouterResponse) {
	if t.next != nil {
		f := t.next
		t.next = nil
		f()
	}
}

func (t *Tag) Then(f func()) {
	t.next = f
}

func (t *Tag) Type() string {
	if _, ok := TypeMap[t.xtype]; !ok {
		return fmt.Sprintf("%#x", t.xtype)
	} else {
		return TypeMap[t.xtype]
	}
}

func (t *Tag) Name() string {
	return string(t.name)
}

func NewTag(name string) *Tag {

	_tag := &Tag{}
	_tag.name = []byte(name)
	_tag.readCount = 1
	_tag.xtype = NULL
	return _tag
}

func (t *Tag) GeneratePath() ([]byte, error) {
	path := make([]byte, 0)

	pathArr := make([]string, 0)

	for _, str := range regexp.MustCompile(`[.[\],]`).Split(t.Name(), -1) {
		if len(str) > 0 {
			pathArr = append(pathArr, str)
		}
	}

	lenPathArr := len(pathArr)
	bitIndex := int(-1)

	memArr := strings.Split(t.Name(), ".")
	lenMemArr := len(memArr)

	isBitIndex := (lenMemArr > 0 && isValInt(memArr[lenMemArr-1]))

	isBitString := (t.xtype == BIT_STRING && isValInt(pathArr[lenPathArr-1]))

	if isBitString && isBitIndex {
		return path, fmt.Errorf("Tag cannot be defined as a BIT_STRING and have a bit index")
	}

	if isBitString {
		bitIndexConvertion, err := strconv.Atoi(pathArr[lenPathArr-1])

		if err != nil {
			return path, fmt.Errorf("Wrong isBitString Convertion")
		}

		bitIndex = bitIndexConvertion % 32
		pathArr[lenPathArr-1] = fmt.Sprint((bitIndexConvertion - bitIndex) / 32)

	} else {
		if isBitIndex {
			bitIndex, err := strconv.Atoi(pathArr[lenPathArr-1])

			if err != nil {
				return path, fmt.Errorf("Wrong bitIndex Convertion")
			}

			if bitIndex < 0 || bitIndex > 31 {
				return path, fmt.Errorf("Tag bit index must be between 0 and 31")
			}

		}
	}

	buf := &bytes.Buffer{}
	for _, p := range pathArr {
		path = epath.BuildRev(p)
	}

	path = buf.Bytes()
	return path, nil
}

func isValInt(val string) bool {
	if _, err := strconv.Atoi(val); err == nil {
		return true
	}
	return false
}

func NewTagWithType(name string, tp DataType) *Tag {
	_tag := &Tag{}
	_tag.name = []byte(name)
	_tag.readCount = 1
	_tag.xtype = tp
	return _tag
}

func (t *Tag) GetValue() interface{} {
	reader := bytes.NewReader(t.value)

	switch t.xtype {
	case NULL:
		return nil
	case SINT:
		result := int8(0)
		lib.ReadByte(reader, &result)
		return result
	case INT:
		result := int16(0)
		lib.ReadByte(reader, &result)
		return result
	case DINT:
		result := int32(0)
		lib.ReadByte(reader, &result)
		return result
	case LINT:
		result := int64(0)
		lib.ReadByte(reader, &result)
		return result
	case REAL:
		result := float32(0)
		lib.ReadByte(reader, &result)
		return result
	case LREAL:
		result := float64(0)
		lib.ReadByte(reader, &result)
		return result
	case STRUCT:
		_tp1 := uint16(0)
		lib.ReadByte(reader, &_tp1)
		if _tp1 == 0xfce {
			t.structType = STRINGAB
			_len := uint32(0)
			lib.ReadByte(reader, &_len)
			buf := make([]byte, _len)
			lib.ReadByte(reader, buf)
			return string(buf)
		} else {
			return t.value
		}
	default:
		return t.value
	}
}

func (t *Tag) SetValue(data interface{}) {
	writer := new(bytes.Buffer)

	switch t.xtype {
	case NULL:
	case SINT:
		result := data.(int8)
		lib.WriteByte(writer, &result)
	case INT:
		result := data.(int16)
		lib.WriteByte(writer, &result)
	case DINT:
		result := data.(int32)
		lib.WriteByte(writer, &result)
	case LINT:
		result := data.(int64)
		lib.WriteByte(writer, &result)
	case REAL:
		result := data.(float32)
		lib.WriteByte(writer, &result)
	case LREAL:
		result := data.(float64)
		lib.WriteByte(writer, &result)
	}

	if writer.Len() > 0 {
		t.value = writer.Bytes()
	}
}
