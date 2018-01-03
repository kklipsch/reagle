package local

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

//Command is the request body structure
type Command struct {
	XMLName xml.Name `xml:"Command"`
	Name    string   `xml:"Name"`
}

//WifiStatusCommand creates the Command request for wifi status
func WifiStatusCommand() Command {
	return Command{Name: "wifi_status"}
}

//DeviceListCommand creates the Command request for device lists
func DeviceListCommand() Command {
	return Command{Name: "device_list"}
}

//PostManagerEndpoint returns the url to the PostManagerEndpoint
func PostManagerEndpoint(config Config) string {
	return fmt.Sprintf("http://%s/cgi-bin/post_manager", config.Location)
}

//PostCommand posts the provided command to the location using the provided client
func PostCommand(ctx context.Context, client *http.Client, config Config, command interface{}) (code int, body []byte, err error) {
	var (
		commandBody []byte
		req         *http.Request
		resp        *http.Response
	)

	commandBody, err = xml.Marshal(command)
	if err != nil {
		return
	}

	req, err = http.NewRequest("POST", PostManagerEndpoint(config), bytes.NewReader(commandBody))
	if err != nil {
		return
	}

	req.SetBasicAuth(config.User, config.Password)

	resp, err = client.Do(req)
	if err != nil {
		return
	}

	code = resp.StatusCode
	body, err = ioutil.ReadAll(resp.Body)
	return
}
