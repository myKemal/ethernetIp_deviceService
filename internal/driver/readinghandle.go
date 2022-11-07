package driver

import (
	config "device-ethernetip-go/internal/config"
	client "device-ethernetip-go/package/ethernetip_client"
	"device-ethernetip-go/package/tag"
	"time"

	sdkModel "github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
)

func (d *Driver) HandleReadCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []sdkModel.CommandRequest) (res []*sdkModel.CommandValue, err error) {
	res = make([]*sdkModel.CommandValue, len(reqs))

	connectionInfo, err := config.CreateConnectionInfo(protocols)
	if err != nil {
		driver.Logger.Errorf("Fail to create read command connection info. err:%v \n", err)
		return res, err
	}

	err = d.lockAddress(connectionInfo)
	if err != nil {
		return res, err
	}
	defer d.unlockAddress(connectionInfo)

	cfg := client.DefaultConfig()
	cfg.Port = uint16(connectionInfo.Port)

	device, _ := client.NewOriginator(connectionInfo.Address, 0, cfg)

	connectionErr := device.Connect()

	if connectionErr != nil {
		return res, connectionErr
	}

	device.OnConnected = func() {

		for _, req := range reqs {

			_tag := tag.NewTag(req.DeviceResourceName)
			device.ReadTagRev(_tag).Then(func() {

				result, err := sdkModel.NewCommandValue(req.DeviceResourceName, req.Type, res)
				if err != nil {
					return
				}

				result.Origin = time.Now().UnixNano() / int64(time.Millisecond)
				result.Value = _tag.GetValue()
				res = append(res, result)
			})
		}

	}

	d.Logger.Infof("Readed: %v", res)

	return res, nil
}
