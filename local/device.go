package local

import (
	"encoding/xml"
	"strconv"
	"time"
)

//DeviceData is the data about a device, it is named Device in some responses and DeviceDetails in others
type DeviceData struct {
	HardwareAddress  string `json:"hardware_address"`
	Manufacturer     string `json:"manufacturer"`
	ModelID          string `xml:"ModelId" json:"model_id"`
	Protocol         string `json:"protocol"`
	LastContact      string `json:"last_contact"`
	ConnectionStatus string `json:"connection_status"`
	NetworkAddress   string `json:"network_address"`
}

//LastContactTime parses the contact time string into a golang time
//TODO: change this so the unmarshaller does it
func (item DeviceData) LastContactTime() (t time.Time, err error) {
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

//Device is sometimes used
type Device struct {
	DeviceData
	XMLName xml.Name `xml:"Device" json:"-"`
}

//DeviceDestails is also used
type DeviceDetails struct {
	DeviceData
	XMLName xml.Name `xml:"DeviceDetails" json:"-"`
}
