package xray

import "github.com/xtls/xray-core/features/extension"

// NotifyNetworkChanged schedules the observatory check; no probe runs on the
// caller's platform callback thread. Serialize lookup against core shutdown.
func NotifyNetworkChanged() {
	coreServerMu.Lock()
	defer coreServerMu.Unlock()
	if coreServer == nil || !coreServer.IsRunning() {
		return
	}
	if observer, ok := coreServer.GetFeature(extension.ObservatoryType()).(interface{ NotifyNetworkChanged() }); ok {
		observer.NotifyNetworkChanged()
	}
}

func NetworkChangeReady() bool {
	coreServerMu.Lock()
	defer coreServerMu.Unlock()
	if coreServer == nil || !coreServer.IsRunning() {
		return false
	}
	if observer, ok := coreServer.GetFeature(extension.ObservatoryType()).(interface{ NetworkChangeReady() bool }); ok {
		return observer.NetworkChangeReady()
	}
	return true
}
