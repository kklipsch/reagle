package local

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

const (
	//LocationEnv is the name of the environment variable that stores the ip/domain name for the eagle device (e.g. 192.168.100.164)
	LocationEnv string = "REAGLE_LOCAL_LOCATION"
	//UserEnv is the name of the environment variable that stores the username for the local api authentication, usually the cloudid of the eagle device
	UserEnv string = "REAGLE_LOCAL_USER"
	//PasswordEnv is the name of the environment variable that stores the password for the local api authentication, usually the install code of the eagle device
	PasswordEnv string = "REAGLE_LOCAL_PASSWORD"

	//DebugRequestEnv will turn on request debugging if it is any value other than empty
	DebugRequestEnv string = "REAGLE_DEBUG_REQUEST"
	//DebugResponseEnv will turn on response debugging if it is any value other than empty
	DebugResponseEnv string = "REAGLE_DEBUG_RESPONSE"

	//ImprovedFirmwareEnv set to yes if your firmware responds with well formed queries to multiplier and divisor queries, set to no if not
	ImprovedFirmwareEnv string = "REAGLE_IMPROVED_FIRMWARE"

	//MeterModelIDEnv is the name of the 'model_id' returned by the device for the smart meter being watched.  defaults to 'electric_meter' if not set
	MeterModelIDEnv string = "REAGLE_MODEL_ID_NAME"
)

//TestConfigOrSkip returns teh Config from the environment variables or skips if any aren't set
func TestConfigOrSkip(t testing.TB) Config {
	config, ok := ConfigFromEnv()
	if !ok {
		t.Skipf("Skipping because one or more of [%v, %v, %v] is not set", LocationEnv, UserEnv, PasswordEnv)
	}

	return config
}

//ConfigFromEnv returns a Config and true using the environment variables or a Config and false if any aren't set
func ConfigFromEnv() (Config, bool) {
	//unless affirmatively set assume that they have the improved firmware with the bug around variables
	improved := strings.ToLower(strings.TrimSpace(os.Getenv(ImprovedFirmwareEnv))) != "false"

	filter := BadResponseVariables
	if improved {
		filter = NoFilter
	}

	config := Config{
		Location: os.Getenv(LocationEnv),
		User:     os.Getenv(UserEnv),
		password: os.Getenv(PasswordEnv),

		ImprovedFirmware: improved,

		Filter: filter,

		ModelIDForMeter: strings.TrimSpace(os.Getenv(MeterModelIDEnv)),

		DebugRequest:  strings.TrimSpace(os.Getenv(DebugRequestEnv)) == "true",
		DebugResponse: strings.TrimSpace(os.Getenv(DebugResponseEnv)) == "true",
	}

	return config, ConfigOK(config)
}

//ConfigOK returns true if the Config can be used
func ConfigOK(config Config) bool {
	return config.Location != "" && config.User != "" && config.password != ""
}

//ValidateConfig returns an error if the Config is not ready for use
func ValidateConfig(c Config) error {
	if !ConfigOK(c) {
		return fmt.Errorf("Must provide %s, %s, %s: (%s, %s, '*')", LocationEnv, UserEnv, PasswordEnv, c.Location, c.User)
	}

	return nil
}

//SetPassword sets the password on the config
func SetPassword(c Config, password string) Config {
	c.password = password
	return c
}

//Config is used to locate/auth the eagle local api
type Config struct {
	Location string `json:"location"`
	User     string `json:"user"`
	password string `json:"-"`

	//Older versions of the firmware respond with invalid xml for multiplier/divisor
	ImprovedFirmware bool `json:"improved_firmware"`

	DebugRequest  bool `json:"debug_request"`
	DebugResponse bool `json:"debug_response"`

	Filter VariableFilter `json:"variable_filter"`

	//what the eagle returns for the model id of the smart meter to watch.  defaults to electric_meter
	ModelIDForMeter string `json:"model_id"`
}

func (c Config) GetModelIDForMeter() string {
	if c.ModelIDForMeter == "" {
		return "electric_meter"
	}

	return c.ModelIDForMeter
}

func (c Config) GetFilter() VariableFilter {
	if c.Filter == nil {
		return NoFilter
	}

	return c.Filter
}
