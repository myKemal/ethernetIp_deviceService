package driver

import (
	"fmt"
	"sync"

	"device-ethernetip-go/internal/config"

	sdkModel "github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
)

var once sync.Once
var driver *Driver

type Driver struct {
	Logger              logger.LoggingClient
	AsyncCh             chan<- *sdkModel.AsyncValues
	mutex               sync.Mutex
	addressMap          map[string]chan bool
	workingAddressCount map[string]int
	stopped             bool
}

const ColorReset = "\033[0m"
const ColorGreen = "\033[32m"
const ColorRed = "\033[64m"

func NewProtocolDriver() sdkModel.ProtocolDriver {
	once.Do(func() {
		driver = new(Driver)

	})
	return driver
}

func (d *Driver) lockAddress(cI *config.ConnectionInfo) error {
	if d.stopped {
		return fmt.Errorf("service unable to handle new request")
	}

	d.mutex.Lock()

	lock, ok := d.addressMap[cI.Address]

	if !ok {
		lock = make(chan bool, 1)
		d.addressMap[cI.Address] = lock
	}

	count, ok := d.workingAddressCount[cI.Address]

	if !ok {
		d.workingAddressCount[cI.Address] = 1

	} else if count >= cI.CClimit {
		d.mutex.Unlock()
		errorMessage := fmt.Sprintf("High-frequency command execution. There are %v commands with the same address in the queue", cI.CClimit)
		d.Logger.Error(errorMessage)
		return fmt.Errorf(errorMessage)
	} else {
		d.workingAddressCount[cI.Address] = count + 1
	}

	d.mutex.Unlock()
	lock <- true

	return nil
}

func (d *Driver) unlockAddress(cI *config.ConnectionInfo) {
	d.mutex.Lock()
	lock := d.addressMap[cI.Address]
	d.workingAddressCount[cI.Address] = d.workingAddressCount[cI.Address] - 1
	d.mutex.Unlock()
	<-lock
}

func (d *Driver) Initialize(lc logger.LoggingClient, asyncCh chan<- *sdkModel.AsyncValues, deviceCh chan<- []sdkModel.DiscoveredDevice) error {

	d.Logger = lc
	d.AsyncCh = asyncCh
	d.addressMap = make(map[string]chan bool)
	d.workingAddressCount = make(map[string]int)
	return nil
}

func (d *Driver) Stop(force bool) error {
	return nil
}

func (d *Driver) AddDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	d.Logger.Info("Device %s is added", deviceName)
	return nil
}

func (d *Driver) UpdateDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	d.Logger.Info("Device %s is updated", deviceName)
	return nil
}

func (d *Driver) RemoveDevice(deviceName string, protocols map[string]models.ProtocolProperties) error {
	d.Logger.Info("Device %s is removed", deviceName)
	return nil
}
