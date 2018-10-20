package local

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWifiStatus(t *testing.T) {
	config := TestConfigOrSkip(t)
	ctx := context.Background()
	api := New(config)

	status, err := api.WifiStatus(ctx)
	require.NoError(t, err)
	log.Printf("%v", status)
}

func TestDeviceList(t *testing.T) {
	config := TestConfigOrSkip(t)
	ctx := context.Background()
	api := New(config)

	items, err := api.DeviceList(ctx)
	require.NoError(t, err)
	assert.True(t, len(items) > 0)
	log.Printf("%v", items)
}

func TestDeviceQueryDivisorMultiplier(t *testing.T) {
	//The nest returns invalid xml if you try to query the Multiplier or Divisor variables as the description has unescaped & in it
	config := TestConfigOrSkip(t)
	ctx := context.Background()
	api := New(config)

	items, err := api.DeviceList(ctx)
	require.NoError(t, err)
	require.True(t, len(items) > 0)

	_, err := api.DeviceQuery(ctx, items[0].HardwareAddress, "zigbee:Multiplier")
	require.Error(t, err)

	_, err := api.DeviceQuery(ctx, items[0].HardwareAddress, "zigbee:Divisor")
	require.Error(t, err)

}

func TestDeviceQuery(t *testing.T) {
	config := TestConfigOrSkip(t)
	ctx := context.Background()
	api := New(config)

	items, err := api.DeviceList(ctx)
	require.NoError(t, err)
	require.True(t, len(items) > 0)

	resp, err := api.DeviceQuery(ctx, items[0].HardwareAddress, "zigbee:InstantaneousDemand", "zigbee:Message")
	require.NoError(t, err)
	require.True(t, len(resp.Components.Component) > 0, fmt.Sprintf("%v", resp))
	require.Equal(t, len(resp.Components.Component[0].Variables.Variable), 2, fmt.Sprintf("%v", resp))
	require.NotEmpty(t, resp.Components.Component[0].Variables.Variable[0].Value, fmt.Sprintf("%v", resp))

	log.Printf("%v", resp)
}
