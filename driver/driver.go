package driver

import (
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/config"
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/driver/asuswrt"
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/driver/definition"
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/driver/generic"
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/driver/ghw"
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/driver/gopsutil"
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/driver/lmsensors"
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/driver/rpi"
)

// Get Stub Driver (May also be used as base driver for new implementations)
func GetStubDriver() definition.Driver {
	return definition.Driver{
		Is64BitCPU:              func() bool { return false },
		GetBoardModel:           func() string { return "" },
		GetBoardVendor:          func() string { return "" },
		GetSoftwareVersion:      func() string { return "" },
		GetHostId:               func() string { return "" },
		GetHostname:             func() string { return "" },
		GetCPUUsage:             func(bool) []float64 { return []float64{-1} },
		GetCPUTemperature:       func() float64 { return -1 },
		GetDiskUsePercent:       func(string) float64 { return -1 },
		GetRAMUsePercent:        func() float64 { return -1 },
		GetSwapUsePercent:       func() float64 { return -1 },
		GetLastBootTimestamp:    func() uint64 { return 0 },
		GetRaspberryPowerStatus: func() string { return "" },
	}
}

// Get dynamically generated driver from all available drivers

func GetDriver(config *config.Configuration) definition.Driver {

	// Initialize Start Driver
	driver := definition.Driver{}

	// All supported drivers (Comment out the drivers that should not be compiled into the binary)
	rpi.UseDriver(&driver, config)
	asuswrt.UseDriver(&driver, config)
	generic.UseDriver(&driver, config)
	gopsutil.UseDriver(&driver, config)
	ghw.UseDriver(&driver, config)
	lmsensors.UseDriver(&driver, config)

	return driver
}
