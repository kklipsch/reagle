package local

import (
	"encoding/xml"
	"strconv"
	"time"
)

//DeviceListCommand creates the Command request for device lists
func DeviceListCommand() Command {
	return Command{Name: "device_list"}
}

//DeviceList is a list of configured devices on the eagle
type DeviceList struct {
	XMLName xml.Name `xml:"DeviceList"`
	Device  []DeviceListItem
}

//DeviceListItem is a single entry in a DeviceList
type DeviceListItem struct {
	XMLName          xml.Name `xml:"Device"`
	HardwareAddress  string
	Manufacturer     string
	ModelID          string `xml:"ModelId"`
	Protocol         string
	LastContact      string
	ConnectionStatus string
	NetworkAddress   string
}

//LastContactTime parses the contact time string into a golang time
//TODO: change this so the unmarshaller does it
func (item DeviceListItem) LastContactTime() (t time.Time, err error) {
	var (
		i int64
	)

	i, err = strconv.ParseInt(item.LastContact, 0, 64)
	if err != nil {
		return
	}

	t = time.Unix(i, 0)
	return
}
