package local

import "encoding/xml"

//DeviceQueryCommand is the request body structure
type DeviceQueryCommand struct {
	Command
	DeviceDetails DeviceDetails
	Components    Components
}

type DeviceQueryResponse struct {
	XMLName       xml.Name `xml:"Device"`
	DeviceDetails DeviceDetails
	Components    Components
}

func NewDeviceQueryCommand(hardwareAddress string, variables ...string) DeviceQueryCommand {
	return DeviceQueryCommand{
		Command:       NewCommand("device_query"),
		DeviceDetails: DeviceDetails{DeviceData: DeviceData{HardwareAddress: hardwareAddress}},
		Components:    NewComponents(NewComponent("Main", variables...)),
	}
}
