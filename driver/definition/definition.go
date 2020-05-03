package definition

type Is64BitCPU func() bool
type GetBoardModel func() string
type GetBoardVendor func() string
type GetSerialNumber func() string
type GetSoftwareVersion func() string
type GetHostname func() string
type GetTotalCPUUsage func() float64
type GetCPUTemperature func() float64
type GetDiskUsePercent func(string) float64
type GetRAMUsePercent func() float64
type GetSwapUsePercent func() float64
type GetNetworkInBytes func(string) uint64
type GetNetworkOutBytes func(string) uint64
type GetLastBootTimestamp func() uint64
type GetRaspberryPowerStatus func() string

type Driver struct {
	Is64BitCPU Is64BitCPU
	GetBoardModel GetBoardModel
	GetBoardVendor GetBoardVendor
	GetSoftwareVersion GetSoftwareVersion
	GetHostId GetSerialNumber
	GetHostname GetHostname
	GetTotalCPUUsage GetTotalCPUUsage
	GetCPUTemperature GetCPUTemperature
	GetDiskUsePercent GetDiskUsePercent
	GetRAMUsePercent GetRAMUsePercent
	GetSwapUsePercent GetSwapUsePercent
	GetNetworkInBytes GetNetworkInBytes
	GetNetworkOutBytes GetNetworkOutBytes
	GetLastBootTimestamp GetLastBootTimestamp
	GetRaspberryPowerStatus GetRaspberryPowerStatus
}
