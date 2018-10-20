package local

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleCommand() {
	command := Command{Name: "wifi_status"}

	output, err := xml.Marshal(command)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	os.Stdout.Write(output)
	// Output:
	// <Command><Name>wifi_status</Name></Command>
}

func TestPostCommand(t *testing.T) {
	config := TestConfigOrSkip(t)
	ctx := context.Background()

	code, body, err := PostCommand(ctx, &http.Client{}, config, NewDeviceListCommand())
	require.NoError(t, err, fmt.Sprintf("%v - %s", code, body))

	assert.Equal(t, http.StatusOK, code)
	log.Printf("%s", body)
}
