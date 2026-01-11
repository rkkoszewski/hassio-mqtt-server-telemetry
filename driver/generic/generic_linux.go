//go:build linux
// +build linux

package generic

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/config"
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/driver/definition"
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/utils"
)

// Get CPU Temperature from File
func getCPUTemperatureFromFile(path string) float64 {
	tempString, err := utils.ReadFileToString(path)
	if err != nil {
		return -1
	}

	temp, err := strconv.Atoi(strings.Trim(tempString, "\n"))
	if err != nil {
		return -1
	}

	return float64(temp) / float64(1000)
}

func buildCPUTemperatureFromFileFunc(path string) func() float64 {
	return func() float64 {
		return getCPUTemperatureFromFile(path)
	}
}

// Get GPU Usage from File
func getGPUUsageFromFile(path string) float64 {
	usageString, err := utils.ReadFileToString(path)
	if err != nil {
		return -1
	}

	usage, err := strconv.ParseFloat(strings.TrimSpace(usageString), 64)
	if err != nil {
		return -1
	}

	return usage
}

func buildGPUUsageFromFileFunc(path string) func() float64 {
	return func() float64 {
		return getGPUUsageFromFile(path)
	}
}

// Use Generic Linux Driver
func UseDriver(driver *definition.Driver, config *config.Configuration) {

	// Get CPU bitness
	if driver.Is64BitCPU == nil {
		driver.Is64BitCPU = Is64BitCPU
	}

	// Board Model
	if driver.GetBoardModel == nil {
		model, err := utils.ReadFileToString("/sys/devices/virtual/dmi/id/board_name")
		if err == nil {
			driver.GetBoardModel = func() string {
				return model
			}
		}
	}

	// Board Vendor
	if driver.GetBoardVendor == nil {
		vendor, err := utils.ReadFileToString("/sys/devices/virtual/dmi/id/board_vendor")
		if err == nil {
			driver.GetBoardVendor = func() string {
				return vendor
			}
		}
	}

	// CPU Temperature
	if driver.GetCPUTemperature == nil {
		files, err := ioutil.ReadDir("/sys/class/thermal")
		if err == nil {
			for _, f := range files {
				if !f.IsDir() {
					if strings.HasPrefix(f.Name(), "thermal_zone") {
						sensor_name, err := utils.ReadFileToString("/sys/class/thermal/" + f.Name() + "/type")
						if err == nil {
							sensor_name = strings.TrimSpace(sensor_name)

							if strings.EqualFold(sensor_name, "x86_pkg_temp") {
								log.Println("Selecting CPU Temperature sensor: x86_pkg_temp")
								driver.GetCPUTemperature = buildCPUTemperatureFromFileFunc("/sys/class/thermal/" + f.Name() + "/temp")
								return
							}

							if strings.EqualFold(sensor_name, "cpu-thermal") {
								log.Println("Selecting CPU Temperature sensor: cpu-thermal")
								driver.GetCPUTemperature = buildCPUTemperatureFromFileFunc("/sys/class/thermal/" + f.Name() + "/temp")
								return
							}

							if strings.EqualFold(sensor_name, "broadcomThermalDrv") {
								log.Println("Selecting CPU Temperature sensor: broadcomThermalDrv")
								driver.GetCPUTemperature = buildCPUTemperatureFromFileFunc("/sys/class/thermal/" + f.Name() + "/temp")
								return
							}
						}
					}
				}
			}
		}

		if utils.CanReadFile("/sys/class/thermal/thermal_zone0/temp") {
			log.Println("Selecting CPU Temperature sensor in themal_zone0")
			driver.GetCPUTemperature = buildCPUTemperatureFromFileFunc("/sys/class/thermal/thermal_zone0/temp")
		} else {
			if utils.CanReadFile("/sys/class/thermal/thermal_zone1/temp") {
				log.Println("Selecting CPU Temperature sensor in themal_zone1")
				driver.GetCPUTemperature = buildCPUTemperatureFromFileFunc("/sys/class/thermal/thermal_zone1/temp")
			}
		}
	}

	// GPU Usage
	if driver.GetGPUUsage == nil {
		matches, err := filepath.Glob("/sys/class/drm/card*/device/gpu_busy_percent")
		if err == nil && len(matches) > 0 {
			// Use the first available GPU
			log.Printf("Selecting GPU Usage sensor: %s", matches[0])
			driver.GetGPUUsage = buildGPUUsageFromFileFunc(matches[0])
		}
	}

}
