// +build linux

package rpi

import (
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/driver/definition"
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/utils"
	"strings"
)

// Get Raspberry Pi Power Status
func GetRaspberryPowerStatus() string {
	status, err := utils.ReadFileToString("/sys/devices/platform/soc/soc:firmware/get_throttled")
	if err != nil{
		return "readfail"
	}

	switch status {
	case "0\n": return "normal"
	case "1000\n": return "undervolt"
	case "2000\n": return "bad_power_supply"
	case "3000\n": return "bad_power_supply"
	case "4000\n": return "bad_power_supply_throttle"
	case "5000\n": return "normal"
	case "8000\n": return "overheat"
	default: return "unknown"
	}
}

// Use Raspberry Pi Specific Driver
func UseDriver(driver *definition.Driver) {

	// Board Model
	if driver.GetBoardModel == nil {
		model, err := utils.ReadFileToString("/proc/device-tree/model")
		if err == nil{
			model = strings.Trim(model, "\u0000")
			driver.GetBoardModel = func() string {
				return model
			}

			// Detect Raspberry Pi Manufacturer
			if strings.Contains(strings.ToUpper(model), "RASPBERRY PI") {
				driver.GetBoardVendor = func() string {
					return "Raspberry Pi Foundation"
				}
			}
		}
	}

	// Raspberry Pi Power Status
	if driver.GetRaspberryPowerStatus == nil {
		if utils.CanReadFile("/sys/devices/platform/soc/soc:firmware/get_throttled") {
			driver.GetRaspberryPowerStatus = GetRaspberryPowerStatus
		}
	}

}