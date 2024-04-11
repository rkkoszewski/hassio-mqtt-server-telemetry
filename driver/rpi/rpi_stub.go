//go:build !linux
// +build !linux

package rpi

import (
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/config"
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/driver/definition"
)

func UseDriver(driver *definition.Driver, config *config.Configuration) {}
