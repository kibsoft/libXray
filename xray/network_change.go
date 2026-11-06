package xray

import (
	"context"
	"github.com/xtls/xray-core/features/extension"
	"github.com/xtls/xray-core/features/outbound"
	"github.com/xtls/xray-core/transport/internet/splithttp"
)

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
	splithttp.ResetConnections()
	if manager, ok := coreServer.GetFeature(outbound.ManagerType()).(outbound.Manager); ok {
		for _, handler := range manager.ListHandlers(context.Background()) {
			if resetter, ok := handler.(interface{ ResetNetworkConnections() }); ok {
				resetter.ResetNetworkConnections()
			}
		}
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
