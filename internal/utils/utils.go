package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/errors"
)

func ValidIPAddress(IP string) bool {
	switch {
	case isIPv4(IP) || isIPv6(IP):
		return true
	default:
		return false
	}
}

func isIPv4(IP string) bool {
	if !strings.Contains(IP, ".") {
		return false
	}

	ss := strings.Split(IP, ".")

	if len(ss) != 4 {
		return false
	}

	for _, s := range ss {
		if !isV4Num(s) {
			return false
		}
	}

	return true
}

func isV4Num(s string) bool {
	if len(s) == 0 {
		return false
	}

	if len(s) > 1 &&
		(s[0] < '1' || '9' < s[0]) {
		return false
	}

	n, err := strconv.Atoi(s)

	return err == nil && 0 <= n && n < 256
}

func isIPv6(IP string) bool {
	if !strings.Contains(IP, ":") {
		return false
	}

	ss := strings.Split(IP, ":")

	if len(ss) != 8 {
		return false
	}

	for _, s := range ss {
		if !isV6Num(s) {
			return false
		}
	}

	return true
}

func isV6Num(s string) bool {
	if len(s) == 0 || len(s) > 4 {
		return false
	}

	if !('0' <= s[0] && s[0] <= '9') &&
		!('a' <= s[0] && s[0] <= 'z') &&
		!('A' <= s[0] && s[0] <= 'Z') {
		return false
	}

	n, err := strconv.ParseInt(s, 16, 64)

	return err == nil && 0 <= n && n < 1<<16
}

func ParseIntValue(properties map[string]string, key string) (int, error) {
	str, ok := properties[key]
	if !ok {
		return 0, fmt.Errorf("protocol config  not exist")
	}

	val, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("fail to parse protocol config")
	}
	return val, nil
}

func ParseUInt16Value(properties map[string]string, key string) (uint16, error) {
	str, ok := properties[key]
	var base = 16
	var size = 16
	if !ok {
		return 0, fmt.Errorf("protocol config  not exist")
	}
	val, err := strconv.ParseUint(str, base, size)
	if err != nil {
		return 0, fmt.Errorf("fail to parse protocol config ")
	}
	return uint16(val), nil
}

func ParseBoolValue(str string) (bool, error) {

	val, err := strconv.ParseBool(str)

	if err != nil {
		return false, fmt.Errorf("fail to parse protocol config")
	}

	return val, nil

}

func CastSwapAttribute(i interface{}) (res bool, err errors.EdgeX) {
	switch v := i.(type) {
	case bool:
		res = v
	default:
		return res, errors.NewCommonEdgeX(errors.KindContractInvalid, "swap attribute should be boolean value", nil)
	}
	return res, nil
}
