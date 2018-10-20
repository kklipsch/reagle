package local

import "encoding/xml"

//NewWifiStatusCommand creates the Command request for wifi status
func NewWifiStatusCommand() Command {
	return Command{Name: "wifi_status"}
}

//WifiStatus is the response from the wifi_status command
type WifiStatus struct {
	XMLName    xml.Name `xml:"WiFiStatus"`
	Enabled    string
	Type       string
	SSID       string
	Encryption string
	Channel    string
	IPAddress  string `xml:"IpAddress"`
}
