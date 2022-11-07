package ethernetip

import (
	"errors"

	"device-ethernetip-go/package/constants"
)

func RequestNop(data []byte) (*Encapsulation, error) {
	if len(data) > 65511 {
		return nil, errors.New(constants.ErrEncapsulation)
	}

	encapsulation := &Encapsulation{}
	encapsulation.Command = CommandNOP

	return encapsulation, nil
}

func HandleNop(encapsulation *Encapsulation) {

}
