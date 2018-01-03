package local

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
)

//New returns an API with a default http client and the provided config
func New(config Config) API {
	return API{&http.Client{}, config}
}

//API wraps up the eagle local api for ease of use
type API struct {
	Client *http.Client
	Config Config
}

//DeviceList returns the configured devices
func (a API) DeviceList(ctx context.Context) ([]DeviceListItem, error) {
	deviceList := DeviceList{}
	err := a.post(ctx, DeviceListCommand(), &deviceList)
	return deviceList.Device, err
}

//WifiStatus returns the wifi status of the eagle 200
func (a API) WifiStatus(ctx context.Context) (WifiStatus, error) {
	status := WifiStatus{}
	err := a.post(ctx, WifiStatusCommand(), &status)
	return status, err
}

func (a API) post(ctx context.Context, command interface{}, result interface{}) error {
	code, body, err := PostCommand(ctx, a.Client, a.Config, command)
	if err != nil {
		return err
	}

	return unmarshal(code, body, result)
}

func unmarshal(code int, body []byte, v interface{}) error {
	err := xml.Unmarshal(body, v)
	if err != nil {
		return &unmarshalError{code, body, err}
	}

	return nil
}

type unmarshalError struct {
	code int
	body []byte
	err  error
}

func (u *unmarshalError) Error() string {
	return fmt.Sprintf("unable to unmarshal %v: %v - %s", u.err, u.code, u.body)
}