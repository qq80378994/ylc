package system

import (
	"syscall"
)

var (
	user32                    = syscall.MustLoadDLL("user32.dll")
	getSystemMetrics          = user32.MustFindProc("GetSystemMetrics")
	screenWidth, screenHeight int
)

func init() {
	screenWidth = int(getSystemMetricsByIndex(0))
	screenHeight = int(getSystemMetricsByIndex(1))
}

func getSystemMetricsByIndex(index int) uintptr {
	ret, _, _ := getSystemMetrics.Call(uintptr(index))
	return ret
}

func GetWidth() int {
	return screenWidth
}
func GetHeight() int {
	return screenHeight
}
