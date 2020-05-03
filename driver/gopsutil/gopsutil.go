package gopsutil

import (
	"github.com/rkkoszewski/hassio-mqtt-server-telemetry/driver/definition"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"os"
	"strings"
	"time"
)

func GetTotalCPUUsage() float64 {
	percent, err := cpu.Percent(0, false)
	if err != nil {
		return -1
	}

	return percent[0]
}

func GetDiskUsePercent(path string) float64 {
	usage, err := disk.Usage(path)
	if err != nil {
		return -1
	}

	return usage.UsedPercent
}

func GetRAMUsePercent() float64 {
	ram, err := mem.VirtualMemory()
	if err != nil {
		return -1
	}

	return ram.UsedPercent
}

func GetSwapUsePercent() float64 {
	ram, err := mem.SwapMemory()
	if err != nil {
		return -1
	}

	return ram.UsedPercent
}

var (
	netIOCountersStatCacheFirst []net.IOCountersStat
	netIOCountersStatTimeFirst time.Time
	netIOCountersStatCacheSecond []net.IOCountersStat
	netIOCountersStatTimeSecond time.Time
)

// Fetch Network Stats Once a Second
func fetchNetworkStatSample(){
	if time.Since(netIOCountersStatTimeFirst).Seconds() < 1 {
		return
	}

	stat, err := net.IOCounters(true)
	if err != nil {
		return
	}

	// Swap Previous Values
	if netIOCountersStatCacheFirst == nil { // This allows to detect the interface when no previous data is available
		netIOCountersStatCacheSecond = stat
		netIOCountersStatTimeSecond = time.Now()
	}else{
		netIOCountersStatCacheSecond = netIOCountersStatCacheFirst
		netIOCountersStatTimeSecond = netIOCountersStatTimeFirst
	}

	// Store Updated Values
	netIOCountersStatCacheFirst = stat
	netIOCountersStatTimeFirst = time.Now()
}

func GetNetworkInBytes(nic string) uint64 {
	fetchNetworkStatSample()

	for _, statNow := range netIOCountersStatCacheFirst {
		if statNow.Name == nic {
			for _, statBefore := range netIOCountersStatCacheSecond {
				if statBefore.Name == nic {
					byteDiff := statNow.BytesRecv - statBefore.BytesRecv
					timeDiff := netIOCountersStatTimeFirst.Sub(netIOCountersStatTimeSecond)
					return uint64(float64(byteDiff) / timeDiff.Seconds())
				}
			}
		}
	}
	return ^uint64(0)
}

func GetNetworkOutBytes(nic string) uint64 {
	fetchNetworkStatSample()

	for _, statNow := range netIOCountersStatCacheFirst {
		if statNow.Name == nic {
			for _, statBefore := range netIOCountersStatCacheSecond {
				if statBefore.Name == nic {
					byteDiff := statNow.BytesSent - statBefore.BytesSent
					timeDiff := netIOCountersStatTimeFirst.Sub(netIOCountersStatTimeSecond)
					return uint64(float64(byteDiff) / timeDiff.Seconds())
				}
			}
		}
	}
	return ^uint64(0)
}

// Use GOPSUtil based Driver
func UseDriver(driver *definition.Driver){

	if driver.GetTotalCPUUsage == nil && GetTotalCPUUsage() != -1 {
		driver.GetTotalCPUUsage = GetTotalCPUUsage
	}

	if driver.GetDiskUsePercent == nil {
		driver.GetDiskUsePercent = GetDiskUsePercent
	}

	if driver.GetRAMUsePercent == nil && GetRAMUsePercent() != -1 {
		driver.GetRAMUsePercent = GetRAMUsePercent
	}

	if driver.GetSwapUsePercent == nil && GetTotalCPUUsage() != -1 {
		driver.GetSwapUsePercent = GetSwapUsePercent
	}

	if driver.GetHostId == nil || driver.GetSoftwareVersion == nil || driver.GetLastBootTimestamp == nil {

		// Get Host Information
		info, err := host.Info()
		if err == nil{

			// Host ID
			if driver.GetHostId == nil {
				hostId := info.HostID
				driver.GetHostId = func() string {
					return hostId
				}
			}

			// Software Version
			if driver.GetSoftwareVersion == nil {
				softwareVersion := info.PlatformVersion + " (" + strings.Title(info.Platform) + ")"
				driver.GetSoftwareVersion = func() string {
					return softwareVersion
				}
			}

			// Boot Timestamp
			if driver.GetLastBootTimestamp == nil {
				bootTimestamp := info.BootTime
				driver.GetLastBootTimestamp = func() uint64 {
					return bootTimestamp
				}
			}
		}
	}

	// Host Name
	if driver.GetHostname == nil {
		hostname, err := os.Hostname()
		if err == nil {
			driver.GetHostname = func() string {
				return hostname
			}
		}
	}

	if driver.GetHostname == nil {
		hostname, err := os.Hostname()
		if err == nil {
			driver.GetHostname = func() string {
				return hostname
			}
		}
	}

	if driver.GetNetworkInBytes == nil {
		driver.GetNetworkInBytes = GetNetworkInBytes
	}

	if driver.GetNetworkOutBytes == nil {
		driver.GetNetworkOutBytes = GetNetworkOutBytes
	}
}