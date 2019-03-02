package client

import (
	"context"
	"testing"
	"time"

	"github.com/kklipsch/reagle/local"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMediateWifiStatus(t *testing.T) {
	t.Skip("TODO: make this work")

	ctx, clean := context.WithTimeout(context.Background(), time.Second)
	defer clean()

	ts := local.StartTestServer(local.ServeWifiStatus(local.WifiStatus{Enabled: "enabled", SSID: "ssid"}))
	defer ts.Close()

	api := local.New(local.Config{Location: ts.URL})
	mediator := newMediator(api, time.Second)

	result, err := mediator.query(ctx, localWifiStatus, nil)
	require.NoError(t, err)

	status, ok := result.(local.WifiStatus)
	require.True(t, ok)

	assert.Equal(t, "enabled", status.Enabled)
	assert.Equal(t, "ssid", status.SSID)
}
