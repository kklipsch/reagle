package local

import (
	"encoding/xml"
)

//DeviceListCommand creates the Command request for device lists
func DeviceListCommand() Command {
	return Command{Name: "device_list"}
}

//DeviceList is a list of configured devices on the eagle
type DeviceList struct {
	XMLName xml.Name `xml:"DeviceList"`
	Device  []Device
}
