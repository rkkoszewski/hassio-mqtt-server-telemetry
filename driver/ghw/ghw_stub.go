//go:build noghw
// +build noghw

package ghw

import (
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/config"
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/driver/definition"
)

func UseDriver(driver *definition.Driver, config *config.Configuration) {}
