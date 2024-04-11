//go:build linux
// +build linux

package lmsensors

import (
	"strconv"
	"strings"

	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/config"
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/driver/definition"
	"github.com/ssimunic/gosensors"
)

func strip(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if b == '.' ||
			b == ',' ||
			('0' <= b && b <= '9') {
			result.WriteByte(b)
		}
	}
	return result.String()
}

func GetCPUTemperature() float64 {
	sensors, err := gosensors.NewFromSystem()

	if err != nil {
		return -1
	}

	for chip := range sensors.Chips {
		// Iterate over entries
		for key, value := range sensors.Chips[chip] {
			// AMD CPU Temperature
			if key == "Tctl" {
				intVal, err := strconv.ParseFloat(strip(value), 64)

				if err != nil {
					return -1
				}

				return intVal
			}
		}
	}

	return -1
}

// Use Asus Router Specific Driver
func UseDriver(driver *definition.Driver, config *config.Configuration) {

	// CPU Temperature
	if driver.GetCPUTemperature == nil && GetCPUTemperature() != -1 {
		driver.GetCPUTemperature = GetCPUTemperature
	}
}
