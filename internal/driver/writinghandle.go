package driver

import (
	config "device-ethernetip-go/internal/config"
	client "device-ethernetip-go/package/ethernetip_client"

	sdkModel "github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
)

func (d *Driver) HandleWriteCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []sdkModel.CommandRequest, params []*sdkModel.CommandValue) error {

	connectionInfo, err := config.CreateConnectionInfo(protocols)
	if err != nil {
		driver.Logger.Errorf("Fail to create read command connection info. err:%v \n", err)
		return err
	}

	err = d.lockAddress(connectionInfo)
	if err != nil {
		return err
	}
	defer d.unlockAddress(connectionInfo)

	return nil
}

func handleWriteCommandRequest(deviceClient *client.EdgeDevice, req sdkModel.CommandRequest, param *sdkModel.CommandValue) error {
	var err error

	driver.Logger.Infof("Write command finished. Cmd:%v \n", req.DeviceResourceName)

	return err
}
