package config

import (
	"fmt"

	"device-ethernetip-go/internal/utils"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
)

type ConnectionInfo struct {
	Address string `json:"Address"`
	CClimit int    `json:ConcurrentCommandLimit"`
	Port    int    `json:"Port"`
	// Connect & Read timeout(seconds)
	TimeOut int `json:"Timeout"`
	// Idle timeout(seconds) to close the connection
	IdleTimeOut int `json:"IdleTimeout"`
}

func CreateConnectionInfo(protocols map[string]models.ProtocolProperties) (info *ConnectionInfo, err error) {
	errorMessage := "unable to create CIP protocol Driver connection, protocol config '%s'  wrong"
	parseErrorMessage := "parse error  '%s'  wrong format"

	protocolEIP, isExist := protocols[Protocol_ethernetip]

	if !isExist {
		return info, fmt.Errorf("unsupported protocols, must be %s ", Protocol_ethernetip)
	}

	address, ok := protocolEIP[Address]
	if !ok {
		return nil, fmt.Errorf(errorMessage, Address)
	}

	if !utils.ValidIPAddress(address) {
		return nil, fmt.Errorf("addres '%s' is not valid IP", Address)
	}

	cCL, err := utils.ParseIntValue(protocolEIP, ConcurrentCommandLimit)
	if err != nil {
		return nil, fmt.Errorf(parseErrorMessage, ConcurrentCommandLimit)
	}

	port, err := utils.ParseIntValue(protocolEIP, Port)
	if err != nil {
		return nil, fmt.Errorf(parseErrorMessage, Port)
	}

	timeOut, err := utils.ParseIntValue(protocolEIP, TimeOut)
	if err != nil {
		return nil, fmt.Errorf(parseErrorMessage, TimeOut)
	}

	idleTimeOut, err := utils.ParseIntValue(protocolEIP, IdleTimeOut)
	if err != nil {
		return nil, fmt.Errorf(parseErrorMessage, IdleTimeOut)
	}

	return &ConnectionInfo{
		Address:     address,
		CClimit:     cCL,
		Port:        port,
		TimeOut:     timeOut,
		IdleTimeOut: idleTimeOut,
	}, nil
}
