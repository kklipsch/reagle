package local

import "encoding/xml"

//DeviceDetailsCommand gets all variables available
type DeviceDetailsCommand struct {
	Command
	DeviceDetails DeviceDetails
}

type DeviceDetailsResponse struct {
	XMLName       xml.Name       `xml:"Device" json:"-"`
	DeviceDetails DeviceDetails  `json:"details"`
	Components    ComponentNames `json:"components"`
}

func NewDeviceDetailsCommand(hardwareAddress string) DeviceDetailsCommand {
	return DeviceDetailsCommand{
		Command:       NewCommand("device_details"),
		DeviceDetails: DeviceDetails{DeviceData: DeviceData{HardwareAddress: hardwareAddress}},
	}
}

//VariablesFromDetailsResponse unwraps a DeviceDetailsResponse to get the list of Variable names
func VariablesFromDetailsResponse(response DeviceDetailsResponse) []string {
	var variables []string
	for _, component := range response.Components.Component {
		for _, variable := range component.Variables.Variable {
			variables = append(variables, variable)
		}
	}

	return variables
}
