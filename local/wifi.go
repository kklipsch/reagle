package local

import "encoding/xml"

//NewWifiStatusCommand creates the Command request for wifi status
func NewWifiStatusCommand() Command {
	return Command{Name: "wifi_status"}
}

//WifiStatus is the response from the wifi_status command
type WifiStatus struct {
	XMLName    xml.Name `xml:"WiFiStatus" json:"-"`
	Enabled    string   `json:"enabled"`
	Type       string   `json:"type"`
	SSID       string   `json:"ssid"`
	Encryption string   `json:"encryption"`
	Channel    string   `json:"channel"`
	IPAddress  string   `xml:"IpAddress" json:"ip_address"`
}
