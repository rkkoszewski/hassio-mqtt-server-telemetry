package generic

import (
	"strconv"
)

// Returns if the CPU bitness based on the binary build
func Is64BitCPU() bool {
	return strconv.IntSize == 64
}