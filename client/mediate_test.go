package client

import (
	"context"
	"testing"
	"time"

	"github.com/kklipsch/reagle/local"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMediateQuery(t *testing.T) {
	for _, tc := range []mediateTest{
		wifiStatusCheck(),
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, clean := context.WithTimeout(context.Background(), time.Second)
			defer clean()

			ts, config := local.StartTestServer(tc.testServer)
			defer ts.Close()

			api := local.New(config)
			mediator := newMediator(api, time.Second)

			result, err := mediator.query(ctx, tc.typ, tc.payload)
			require.NoError(t, err)

			tc.check(t, result)
		})
	}
}

type mediateTest struct {
	name       string
	testServer local.TestServerPayload
	typ        requestType
	payload    interface{}
	check      func(*testing.T, interface{})
}

func wifiStatusCheck() mediateTest {
	return mediateTest{
		name:       "wifi_status",
		typ:        localWifiStatus,
		testServer: local.ServeWifiStatus(local.WifiStatus{Enabled: "enabled", SSID: "ssid"}),
		check: func(t *testing.T, result interface{}) {
			status, ok := result.(local.WifiStatus)
			require.True(t, ok)

			assert.Equal(t, "enabled", status.Enabled)
			assert.Equal(t, "ssid", status.SSID)
		},
	}
}
