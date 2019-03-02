package local

import "encoding/xml"

//DeviceQueryCommand allows you to query variables
type DeviceQueryCommand struct {
	Command
	DeviceDetails DeviceDetails
	Components    Components
}

type DeviceQueryResponse struct {
	XMLName       xml.Name      `xml:"Device" json:"-"`
	DeviceDetails DeviceDetails `json:"details"`
	Components    Components    `json:"components"`
}

func NewDeviceQueryCommand(hardwareAddress string, variables ...string) DeviceQueryCommand {
	return DeviceQueryCommand{
		Command:       NewCommand("device_query"),
		DeviceDetails: DeviceDetails{DeviceData: DeviceData{HardwareAddress: hardwareAddress}},
		Components:    NewComponents(NewComponent("Main", variables...)),
	}
}
