//go:build linux
// +build linux

package asuswrt

import (
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/driver/definition"
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/utils"
	"os/exec"
	"regexp"
	"strconv"
)

// Use Asus Router Specific Driver
func UseDriver(driver *definition.Driver) {

	// Check if it's an Asus Router
	if utils.FileExists("/rom/.asusrouter") {

		// Board Model
		if driver.GetBoardModel == nil && utils.CommandExists("nvram") {
			modelBytes, err := exec.Command("nvram", "get", "model").CombinedOutput()
			if err == nil {
				model := string(modelBytes)
				driver.GetBoardModel = func() string {
					return model
				}
			}
		}

		// Board Model
		if driver.GetBoardVendor == nil {
			driver.GetBoardVendor = func() string {
				return "Asus"
			}
		}

		// Firmware Version
		if driver.GetSoftwareVersion == nil && utils.CommandExists("nvram") {
			versionBytes, err := exec.Command("nvram", "get", "buildno").CombinedOutput()
			if err == nil {
				version := string(versionBytes)
				driver.GetSoftwareVersion = func() string {
					return version
				}
			}
		}

		// CPU Temperature
		if driver.GetCPUTemperature == nil && utils.CanReadFile("/proc/dmu/temperature") {
			reg := regexp.MustCompile("[-]?\\d[\\d,]*[\\.]?[\\d{2}]*")
			driver.GetCPUTemperature = func() float64 {
				str, err := utils.ReadFileToString("/proc/dmu/temperature")
				if err != nil {
					return -1
				}

				str = reg.FindString(str)

				temp, err := strconv.ParseFloat(str, 10)
				if err != nil {
					return -2
				}

				return temp
			}
		}

		// Get Host ID
		if driver.GetHostId == nil && utils.CommandExists("nvram") {
			hostIDBytes, err := exec.Command("nvram", "get", "serial_no").CombinedOutput()
			if err == nil {
				hostID := string(hostIDBytes)
				driver.GetHostId = func() string {
					return hostID
				}
			}
		}
	}
}
