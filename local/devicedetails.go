package local

import "encoding/xml"

//DeviceDetailsCommand gets all variables available
type DeviceDetailsCommand struct {
	Command
	DeviceDetails DeviceDetails
}

type DeviceDetailsResponse struct {
	XMLName       xml.Name `xml:"Device"`
	DeviceDetails DeviceDetails
	Components    ComponentNames
}

func NewDeviceDetailsCommand(hardwareAddress string) DeviceDetailsCommand {
	return DeviceDetailsCommand{
		Command:       NewCommand("device_details"),
		DeviceDetails: DeviceDetails{DeviceData: DeviceData{HardwareAddress: hardwareAddress}},
	}
}
