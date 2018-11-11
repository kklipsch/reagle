package local

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

//New returns an API with a default http client and the provided config
func New(config Config) API {
	return API{
		Client: &http.Client{},
		Config: config,
	}
}

//API wraps up the eagle local api for ease of use
type API struct {
	Client *http.Client
	Config Config
}

//GetMeterHardwareAddress returns the hardware address of the smart meter, using the name provided in the config
func (a API) GetMeterHardwareAddress(ctx context.Context) (string, error) {
	devices, err := a.DeviceList(ctx)
	if err != nil {
		return "", err
	}

	search := a.Config.GetModelIDForMeter()
	var models []string
	for _, device := range devices {
		if device.ModelID == search {
			return device.HardwareAddress, nil
		}

		models = append(models, device.ModelID)
	}

	return "", fmt.Errorf("no %v found in device list: %v", search, models)
}

//DeviceDetails returns the available variables
func (a API) DeviceDetails(ctx context.Context, hardwareAddress string) (DeviceDetailsResponse, error) {
	deviceResponse := DeviceDetailsResponse{}
	err := a.post(ctx, NewDeviceDetailsCommand(hardwareAddress), &deviceResponse)
	return deviceResponse, err
}

//DeviceQuery returns the queried variable
func (a API) DeviceQuery(ctx context.Context, hardwareAddress string, variables ...string) (DeviceQueryResponse, error) {
	deviceResponse := DeviceQueryResponse{}

	var toquery []string
	filter := a.Config.GetFilter()
	for _, v := range variables {
		if filter.Exclude(v) {
			log.Printf("Excluding %v due to filter: %v", v, filter)
			continue
		}

		toquery = append(toquery, v)
	}

	if len(toquery) < 1 {
		return deviceResponse, fmt.Errorf("Post filter (%v) no variables were available: %v", filter, variables)
	}

	err := a.post(ctx, NewDeviceQueryCommand(hardwareAddress, toquery...), &deviceResponse)
	return deviceResponse, err
}

//DeviceList returns the configured devices
func (a API) DeviceList(ctx context.Context) ([]Device, error) {
	deviceList := DeviceList{}
	err := a.post(ctx, NewDeviceListCommand(), &deviceList)
	return deviceList.Device, err
}

//WifiStatus returns the wifi status of the eagle 200
func (a API) WifiStatus(ctx context.Context) (WifiStatus, error) {
	status := WifiStatus{}
	err := a.post(ctx, NewWifiStatusCommand(), &status)
	return status, err
}

func (a API) post(ctx context.Context, command interface{}, result interface{}) error {
	code, body, err := PostCommand(ctx, a.Client, a.Config, command)
	if err != nil {
		return fmt.Errorf("%v %v\n %s", code, err, body)
	}

	if a.Config.DebugResponse {
		log.Printf("%v - %s", code, body)
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
