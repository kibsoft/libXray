package libXray

import "github.com/xtls/libxray/xray"

// NotifyNetworkChanged starts a fresh physical-network canary check.
func NotifyNetworkChanged() { xray.NotifyNetworkChanged() }

// NetworkChangeReady reports whether the latest physical-network check finished.
func NetworkChangeReady() bool { return xray.NetworkChangeReady() }
