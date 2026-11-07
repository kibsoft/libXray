package libXray

import (
	"encoding/json"
	"testing"
)

func TestNetworkNotificationsInvokeV3Idle(t *testing.T) {
	for _, method := range []string{"notifyNetworkChanged", "networkChangeReady"} {
		raw := Invoke(`{"apiVersion":3,"method":"` + method + `"}`)
		var response struct {
			Success bool `json:"success"`
			Data    struct {
				Ready bool `json:"ready"`
			} `json:"data"`
		}
		if err := json.Unmarshal([]byte(raw), &response); err != nil {
			t.Fatal(err)
		}
		if !response.Success {
			t.Fatalf("%s failed: %s", method, raw)
		}
		if method == "networkChangeReady" && response.Data.Ready {
			t.Fatal("stopped core must not be ready")
		}
	}
}
