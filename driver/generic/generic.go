//go:build !linux
// +build !linux

package generic

import (
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/config"
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/driver/definition"
)

// Use Generic Driver
func UseDriver(driver *definition.Driver, config *config.Configuration) {

	// Get CPU bitness
	if driver.Is64BitCPU == nil {
		driver.Is64BitCPU = Is64BitCPU
	}

}
