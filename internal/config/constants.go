package config

import (
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
)

var ValueTypeBitCountMap = map[string]uint16{
	common.ValueTypeInt16: 16,
	common.ValueTypeInt32: 32,
	common.ValueTypeInt64: 64,

	common.ValueTypeUint16: 16,
	common.ValueTypeUint32: 32,
	common.ValueTypeUint64: 64,

	common.ValueTypeFloat32: 32,
	common.ValueTypeFloat64: 64,

	common.ValueTypeBool:   1,
	common.ValueTypeString: 16,
}
