package utils

import "strings"

const (
	Byte = 0
	KB = 1
	MB = 2
	GB = 3
	TB = 4
	KB_DIV = float64(1<<10)
	MB_DIV = float64(1<<20)
	GB_DIV = float64(1<<30)
	TB_DIV = float64(1<<40)
)

// Return String Measure to Measure ID
func ConvertNetworkUnitOfMeasureToID(measure string) uint8 {
	measure = strings.ToUpper(measure)
	switch measure {
	default:
		fallthrough
	case "BYTE":
		return 0
	case "KILOBYTE":
		fallthrough
	case "KB":
		return 1
	case "MEGABYTE":
		fallthrough
	case "MB":
		return 2
	case "GIGABYTE":
		fallthrough
	case "GB":
		return 3
	case "TERABYTE":
		fallthrough
	case "TB":
		return 4
	}
}

// Return Network Unit of Measure
func GetNetworkUnitOfMeasure(measureId uint8) string {
	switch measureId {
	default:
		fallthrough
	case Byte:
		return "Bytes/s"
	case KB:
		return "KB/s"
	case MB:
		return "MB/s"
	case GB:
		return "GB/s"
	case TB:
		return "TB/s"
	}
}

// Convert Unit of Measure
func ConvertToUnitOfMeasure(value float64, measureId uint8) float64 {
	switch measureId {
	default:
		fallthrough
	case Byte:
		return value
	case KB:
		return value / KB_DIV
	case MB:
		return value / MB_DIV
	case GB:
		return value / GB_DIV
	case TB:
		return value / TB_DIV
	}
}