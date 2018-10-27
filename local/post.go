package local

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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

	endpoint := PostManagerEndpoint(config)

	if config.DebugRequest {
		log.Printf("%s\n%s", endpoint, commandBody)
	}

	req, err = http.NewRequest("POST", endpoint, bytes.NewReader(commandBody))
	if err != nil {
		return
	}

	req.SetBasicAuth(config.User, config.password)

	resp, err = client.Do(req)
	if err != nil {
		return
	}

	code = resp.StatusCode
	body, err = ioutil.ReadAll(resp.Body)
	return
}
