package local

import (
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
	config := Config{
		Location: os.Getenv(LocationEnv),
		User:     os.Getenv(UserEnv),
		password: os.Getenv(PasswordEnv),

		DebugRequest:  strings.TrimSpace(os.Getenv(DebugRequestEnv)) != "",
		DebugResponse: strings.TrimSpace(os.Getenv(DebugResponseEnv)) != "",
	}

	ok := config.Location != "" && config.User != "" && config.password != ""

	return config, ok
}

//Config is used to locate/auth the eagle local api
type Config struct {
	Location string
	User     string
	password string

	DebugRequest  bool
	DebugResponse bool
}
