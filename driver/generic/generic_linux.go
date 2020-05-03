// +build linux

package generic

import (
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/driver/definition"
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/utils"
	"math"
	"strconv"
	"strings"
)

// Get CPU Temperature from File
func getCPUTemperatureFromFile(path string) float64 {
	tempString, err := utils.ReadFileToString(path)
	if err != nil{
		return -1
	}

	temp, err := strconv.Atoi(strings.Trim(tempString, "\n"))
	if err != nil{
		return -1
	}

	return math.Round((float64(temp) / float64(1000))*100)/100
}

// Use Generic Linux Driver
func UseDriver(driver *definition.Driver) {

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
		if err == nil{
			driver.GetBoardVendor = func() string {
				return vendor
			}
		}
	}

	// CPU Temperature
	if driver.GetCPUTemperature == nil {
		if utils.CanReadFile("/sys/class/thermal/thermal_zone0/temp") {
			driver.GetCPUTemperature = func() float64 {
				return getCPUTemperatureFromFile("/sys/class/thermal/thermal_zone0/temp")
			}
		}else{
			if utils.CanReadFile("/sys/class/thermal/thermal_zone1/temp") {
				driver.GetCPUTemperature = func() float64 {
					return getCPUTemperatureFromFile("/sys/class/thermal/thermal_zone1/temp")
				}
			}
		}
	}

}