package local

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
)

//ServeDeviceList returns a payload for testing DeviceList commands
func ServeDeviceList(list []Device) TestServerPayload {
	return TestServerPayload{DeviceList: list}
}

//ServeDeviceDetails returns a payload for testing DeviceDetails commands
func ServeDeviceDetails(details DeviceDetailsResponse) TestServerPayload {
	return TestServerPayload{DeviceDetails: &details}
}

//ServeDeviceQuery returns a payload for testing DeviceQuery commands
func ServeDeviceQuery(query DeviceQueryResponse) TestServerPayload {
	return TestServerPayload{DeviceQuery: &query}
}

//ServeWifiStatus returns a payload for testing WifiStatus commands
func ServeWifiStatus(status WifiStatus) TestServerPayload {
	return TestServerPayload{WifiStatus: &status}
}

//TestServerPayload is the responses the httptest.Server should respond with
type TestServerPayload struct {
	DeviceList    []Device
	DeviceDetails *DeviceDetailsResponse
	DeviceQuery   *DeviceQueryResponse
	WifiStatus    *WifiStatus
}

//StartTestServer returns an httptest.Server that responds similar to the Eagle
func StartTestServer(payload TestServerPayload) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(testServer(payload)))
}

func testServer(payload TestServerPayload) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/post_manager" {
			http.Error(w, fmt.Sprintf("unknown path: %v", r.URL.Path), http.StatusNotFound)
			return
		}

		if r.Method != "POST" {
			http.Error(w, "must be post", http.StatusMethodNotAllowed)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		command := Command{}
		err = xml.Unmarshal(body, &command)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var response interface{}
		switch command.Name {
		case "device_list":
			response = payload.DeviceList
		case "device_details":
			response = payload.DeviceDetails
		case "device_query":
			response = payload.DeviceQuery
		case "wifi_status":
			response = payload.WifiStatus
		default:
			http.Error(w, fmt.Sprintf("unknown command name: %v", command.Name), http.StatusBadRequest)
			return
		}

		if response == nil {
			http.Error(w, fmt.Sprintf("no payload for command: %v", command.Name), http.StatusInternalServerError)
			return
		}

		b, err := xml.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(b)
		if err != nil {
			log.Fatalf("could not write: %v", err)
		}
		return
	}
}
