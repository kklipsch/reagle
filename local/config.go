package local

import (
	"os"
	"testing"
)

const (
	//LocationEnv is the name of the environment variable that stores the ip/domain name for the eagle device (e.g. 192.168.100.164)
	LocationEnv string = "REAGLE_LOCAL_LOCATION"
	//UserEnv is the name of the environment variable that stores the username for the local api authentication, usually the cloudid of the eagle device
	UserEnv string = "REAGLE_LOCAL_USER"
	//PasswordEnv is the name of the environment variable that stores the password for the local api authentication, usually the install code of the eagle device
	PasswordEnv string = "REAGLE_LOCAL_PASSWORD"
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
		Password: os.Getenv(PasswordEnv),
	}

	ok := config.Location != "" && config.User != "" && config.Password != ""

	return config, ok
}

//Config is used to locate/auth the eagle local api
type Config struct {
	Location string
	User     string
	Password string
}
