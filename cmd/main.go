package main

import (
	device_ethernetip_go "device-ethernetip-go"

	"device-ethernetip-go/internal/driver"

	"github.com/edgexfoundry/device-sdk-go/v2/pkg/startup"
)

const (
	serviceName string = "device-ethernetip-go"
)

func main() {

	dp := driver.NewProtocolDriver()

	startup.Bootstrap(serviceName, device_ethernetip_go.Version, dp)
}
