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

	_, err := api.WifiStatus(ctx)
	require.NoError(t, err)
}

func TestDeviceList(t *testing.T) {
	config := TestConfigOrSkip(t)
	ctx := context.Background()
	api := New(config)

	items, err := api.DeviceList(ctx)
	require.NoError(t, err)
	assert.True(t, len(items) > 0)
}

func TestDeviceQueryDivisorMultiplier(t *testing.T) {
	config := TestConfigOrSkip(t)
	config.Filter = NoFilter //want divsor/multiplier in no matter what in this test

	ctx := context.Background()

	api := New(config)

	items, err := api.DeviceList(ctx)
	require.NoError(t, err)
	require.True(t, len(items) > 0)

	_, err = api.DeviceQuery(ctx, items[0].HardwareAddress, "zigbee:Multiplier")
	if config.ImprovedFirmware {
		assert.NoError(t, err)
	} else {
		assert.Error(t, err)
	}

	_, err = api.DeviceQuery(ctx, items[0].HardwareAddress, "zigbee:Divisor")
	if config.ImprovedFirmware {
		assert.NoError(t, err)
	} else {
		assert.Error(t, err)
	}
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
	require.True(t, len(resp.Components.Component[0].Variables.Variable) > 0, fmt.Sprintf("%v", resp))
	require.NotEmpty(t, resp.Components.Component[0].Variables.Variable[0].Value, fmt.Sprintf("%v", resp))
}

func TestDeviceDetails(t *testing.T) {
	config := TestConfigOrSkip(t)
	ctx := context.Background()
	api := New(config)

	items, err := api.DeviceList(ctx)
	require.NoError(t, err)
	require.True(t, len(items) > 0)

	resp, err := api.DeviceDetails(ctx, items[0].HardwareAddress)
	require.NoError(t, err)
	require.True(t, len(resp.Components.Component) > 0, fmt.Sprintf("%v", resp))
	require.True(t, len(resp.Components.Component[0].Variables.Variable) > 0, fmt.Sprintf("%v", resp))
}

func TestVariables(t *testing.T) {
	if testing.Short() {
		t.Skip("Variable test takes a long time")
	}

	config := TestConfigOrSkip(t)
	ctx := context.Background()
	api := New(config)

	items, err := api.DeviceList(ctx)
	require.NoError(t, err)
	require.True(t, len(items) > 0)

	hardwareAddress := items[0].HardwareAddress
	resp, err := api.DeviceDetails(ctx, hardwareAddress)
	require.NoError(t, err)

	filter := config.GetFilter()
	for _, component := range resp.Components.Component {
		for _, variable := range component.Variables.Variable {
			if filter.Exclude(variable) {
				log.Printf("skipping %s due to variable filter: %v", variable, filter)
				continue
			}

			_, err := api.DeviceQuery(ctx, hardwareAddress, variable)
			require.NoError(t, err)
		}
	}

}
