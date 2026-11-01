package xray

import (
	"context"
	"os"
	"runtime/debug"
	"strconv"
	"sync"

	"github.com/xtls/libxray/memory"
	"github.com/xtls/xray-core/common/cmdarg"
	"github.com/xtls/xray-core/common/platform"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/features/extension"
	"github.com/xtls/xray-core/features/outbound"
	_ "github.com/xtls/xray-core/main/distro/all"
	"github.com/xtls/xray-core/transport/internet/splithttp"
)

var (
	coreServer   *core.Instance
	coreServerMu sync.Mutex
)

// NetworkChangeReady is called under the embedding application's lifecycle lock.
func NetworkChangeReady() bool {
	if coreServer == nil {
		return false
	}
	if observer, ok := coreServer.GetFeature(extension.ObservatoryType()).(interface{ NetworkChangeReady() bool }); ok {
		return observer.NetworkChangeReady()
	}
	return true
}

// NotifyNetworkChanged is called under the embedding application's lifecycle lock.
func NotifyNetworkChanged() {
	if coreServer == nil {
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

func StartXray(configPath string) (*core.Instance, error) {
	file := cmdarg.Arg{configPath}
	config, err := core.LoadConfig("json", file)
	if err != nil {
		return nil, err
	}

	server, err := core.New(config)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func StartXrayFromJSON(configJSON string) (*core.Instance, error) {
	// Convert JSON string to bytes
	configBytes := []byte(configJSON)

	// Use core.StartInstance which can load configuration directly from bytes
	server, err := core.StartInstance("json", configBytes)
	if err != nil {
		return nil, err
	}

	return server, nil
}

// SetTunFd sets the TUN file descriptor.
// Call this BEFORE RunXray/RunXrayFromJSON.
func SetTunFd(fd int32) {
	os.Setenv(platform.TunFdKey, strconv.Itoa(int(fd)))
}

func InitEnv(datDir string) {
	os.Setenv(platform.AssetLocation, datDir)
	os.Setenv(platform.CertLocation, datDir)
}

// Run Xray instance.
// datDir means the dir which geosite.dat and geoip.dat are in.
// configPath means the config.json file path.
func RunXray(datDir, configPath string) (err error) {
	coreServerMu.Lock()
	defer coreServerMu.Unlock()
	InitEnv(datDir)
	memory.InitForceFree()
	coreServer, err = StartXray(configPath)
	if err != nil {
		return
	}

	if err = coreServer.Start(); err != nil {
		return
	}

	debug.FreeOSMemory()
	return nil
}

// Run Xray instance with JSON configuration string.
// datDir means the dir which geosite.dat and geoip.dat are in.
// configJSON means the JSON configuration string.
func RunXrayFromJSON(datDir, configJSON string) (err error) {
	coreServerMu.Lock()
	defer coreServerMu.Unlock()
	InitEnv(datDir)
	memory.InitForceFree()
	coreServer, err = StartXrayFromJSON(configJSON)
	if err != nil {
		return
	}

	debug.FreeOSMemory()
	return nil
}

// Get Xray State
func GetXrayState() bool {
	coreServerMu.Lock()
	defer coreServerMu.Unlock()
	return coreServer != nil && coreServer.IsRunning()
}

// Stop Xray instance.
func StopXray() error {
	coreServerMu.Lock()
	defer coreServerMu.Unlock()
	if coreServer != nil {
		err := coreServer.Close()
		coreServer = nil
		if err != nil {
			return err
		}
	}
	return nil
}

// Xray's version
func XrayVersion() string {
	return core.Version()
}
