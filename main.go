package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/config"
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/driver"
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/hassio"
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/utils"
)

func main() {
	log.Println("Starting Home Assistant MQTT Server Telemetry Service")

	// Read Configuration
	configFile := "configuration.yaml"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	config, err := config.Read(configFile)
	if err != nil {
		log.Fatal("ERROR: ", err.Error())
	}

	// Startup Delay
	if config.Advanced.StartDelay > 0 {
		log.Printf("Startup delay is configured. Sleeping for %v seconds.", config.Advanced.StartDelay)
		time.Sleep(time.Duration(config.Advanced.StartDelay) * time.Second)
	}

	// Setting Variables
	name := config.Name
	updateSleepTimeout := 15 * time.Second
	if config.TelemetryUpdateDelay >= 1 {
		updateSleepTimeout = time.Duration(config.TelemetryUpdateDelay) * time.Second
	}

	//nic := "eth0"
	//diskUsageDrivePath := "/mnt/usbhd0"

	// Initialize Driver
	driver := driver.GetDriver(&config)

	// Gather basic information
	hostname := "Undefined"
	if driver.GetHostname != nil {
		hostname = driver.GetHostname()
	}

	boardModel := "Unknown"
	if driver.GetBoardModel != nil {
		boardModel = driver.GetBoardModel()
	}

	boardVendor := "Unknown"
	if driver.GetBoardVendor != nil {
		boardVendor = driver.GetBoardVendor()
	}

	softwareVersion := "Unknown"
	if driver.GetSoftwareVersion != nil {
		softwareVersion = driver.GetSoftwareVersion()
	}

	hostId := "UniqueDeviceIDMissing"
	if config.Advanced.DeviceID != "auto" {
		hostId = config.Advanced.DeviceID
	} else {
		if driver.GetHostId != nil {
			hostId = driver.GetHostId()
		}
	}

	// Initialize MQTT Client
	client := hassio.NewClient(config.Broker, config.Id, config.User, config.Password)

	// Create Device
	device := client.AddDevice(
		hostname,
		boardModel,
		boardVendor,
		softwareVersion,
		[]string{hostId})

	// Create Device Sensors

	// CPU Usage
	if driver.GetCPUUsage != nil && config.Cpu.Usage.Enabled {
		cpuUsageIcon := "mdi:cpu-32-bit"
		if driver.Is64BitCPU != nil && driver.Is64BitCPU() {
			cpuUsageIcon = "mdi:cpu-64-bit"
		}

		if config.Cpu.Usage.Total {
			device.AddSensor(fmt.Sprintf("%s Total CPU Usage", name),
				"cpu_usage_total",
				nil,
				"%",
				"",
				cpuUsageIcon,
				func() interface{} {
					return utils.ValuePrecision(driver.GetCPUUsage(true)[0], config.Cpu.Usage.Decimal)
				})
		} else {
			// Total CPU Usage
			cpuUsage := driver.GetCPUUsage(false)
			device.AddSensor(fmt.Sprintf("%s Total CPU Usage", name),
				"cpu_usage_total",
				nil,
				"%",
				"",
				cpuUsageIcon,
				func() interface{} {
					// Update the CPU usage values
					cpuUsage = driver.GetCPUUsage(false)

					// Calculate the total CPU usage
					totalUsage := 0.0
					for _, coreUsage := range cpuUsage {
						totalUsage += coreUsage
					}

					return utils.ValuePrecision(totalUsage/float64(len(cpuUsage)), config.Cpu.Usage.Decimal)
				})

			// Per Core CPU Usage
			for i := range cpuUsage {
				currentCore := i
				device.AddSensor(fmt.Sprintf("%s CPU Core %d Usage", name, currentCore),
					fmt.Sprintf("cpu_usage_core_%d", currentCore),
					nil,
					"%",
					"",
					cpuUsageIcon,
					func() interface{} {
						/*  Superfluous, as the CPU usage is updated in the total CPU usage sensor
						// If it's the first index update the CPU usage values
						if currentCore == 0 {
							cpuUsage = driver.GetCPUUsage(false)
						}
						*/
						return utils.ValuePrecision(cpuUsage[currentCore], config.Cpu.Usage.Decimal)
					})
			}
		}
	}

	// CPU Temperature
	if driver.GetCPUTemperature != nil && config.Cpu.Temperature.Enabled {
		device.AddSensor(fmt.Sprintf("%s CPU Temperature", name),
			"cpu_temperature",
			"temperature",
			"°C",
			"",
			"mdi:thermometer",
			func() interface{} {
				return utils.ValuePrecision(driver.GetCPUTemperature(), config.Cpu.Temperature.Decimal)
			})
	}

	// RAM Use
	if driver.GetRAMUsePercent != nil && config.Ram.Enabled {
		device.AddSensor(fmt.Sprintf("%s RAM Use", name),
			"ram_use",
			nil,
			"%",
			"",
			"mdi:memory",
			func() interface{} {
				return utils.ValuePrecision(driver.GetRAMUsePercent(), config.Ram.Decimal)
			})
	}

	// SWAP Use
	if driver.GetSwapUsePercent != nil && config.Swap.Enabled {
		device.AddSensor(fmt.Sprintf("%s SWAP Use", name),
			"swap_use",
			nil,
			"%",
			"",
			"mdi:harddisk",
			func() interface{} {
				return utils.ValuePrecision(driver.GetSwapUsePercent(), config.Swap.Decimal)
			})
	}

	// Disk Use
	if driver.GetDiskUsePercent != nil {
		for _, drive := range config.Storage {
			// Drive ID
			driveId := utils.StripSpecialChars(drive.Drive, true)
			device.AddSensor(fmt.Sprintf("%s Disk Use %s", name, drive.Drive),
				fmt.Sprintf("disk_use_%s", driveId),
				nil,
				"%",
				"",
				"mdi:micro-sd",
				func() interface{} {
					return utils.ValuePrecision(driver.GetDiskUsePercent(drive.Drive), drive.Decimal)
				})
		}
	}

	// Networking
	for _, network := range config.Network {
		// Network ID
		netId := utils.StripSpecialChars(network.Interface, true)
		measureId := utils.ConvertNetworkUnitOfMeasureToID(network.Bitrate)

		// Network In
		if driver.GetNetworkInBytes != nil && network.Ingress && driver.GetNetworkInBytes(network.Interface) != ^uint64(0) {
			device.AddSensor(fmt.Sprintf("%s %s Network In", name, network.Interface),
				fmt.Sprintf("network_in_%s", netId),
				nil,
				utils.GetNetworkUnitOfMeasure(measureId),
				"",
				"mdi:download",
				func() interface{} {
					bytes := driver.GetNetworkInBytes(network.Interface)
					if bytes == 0 {
						return 0
					}
					return utils.ValuePrecision(
						utils.ConvertToUnitOfMeasure(float64(bytes), measureId),
						network.Decimal)
				})
		}

		// Network Out
		if driver.GetNetworkOutBytes != nil && network.Egress && driver.GetNetworkOutBytes(network.Interface) != ^uint64(0) {
			device.AddSensor(fmt.Sprintf("%s %s Network Out", name, network.Interface),
				fmt.Sprintf("network_out_%s", netId),
				nil,
				utils.GetNetworkUnitOfMeasure(measureId),
				"",
				"mdi:upload",
				func() interface{} {
					bytes := driver.GetNetworkOutBytes(network.Interface)
					if bytes == 0 {
						return 0
					}
					return utils.ValuePrecision(
						utils.ConvertToUnitOfMeasure(float64(bytes), measureId),
						network.Decimal)
				})
		}
	}

	if driver.GetRaspberryPowerStatus != nil && config.Rpi.PowerStatus {
		device.AddSensor(fmt.Sprintf("%s Power Status", name),
			"power_status",
			nil,
			nil,
			"",
			"mdi:power-plug",
			func() interface{} {
				return driver.GetRaspberryPowerStatus()
			})
	}

	if driver.GetLastBootTimestamp != nil {
		device.AddSensor(fmt.Sprintf("%s Last Boot", name),
			"last_boot",
			"timestamp",
			nil,
			"",
			"mdi:clock",
			func() interface{} {
				tm := time.Unix(int64(driver.GetLastBootTimestamp()), 0)
				return tm.UTC()
			})
	}

	// Connect to MQTT Broker
	client.Connect()

	// Handle SigTerm
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-exit
		client.Disconnect()
		os.Exit(0)
	}()

	// Sensor Update Loop
	for {
		// Update Sensor Values
		device.UpdateSensorValues()

		// Submit Changes
		device.SubmitChanges()

		// Sleep
		time.Sleep(updateSleepTimeout)
	}
}
