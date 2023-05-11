package util

import (
	"strconv"

	"github.com/go-vgo/robotgo"
)

func MousePress(s string) {
	button := getMouseClick(s)
	robotgo.MouseToggle("down", button)
}

func MouseRelease(s string) {
	button := getMouseClick(s)
	robotgo.MouseToggle("up", button)
}

func MouseMove(s string) {
	x, y := getMousePos(s)
	robotgo.MoveMouseSmooth(x, y)
}

func getMouseClick(s string) uint32 {
	button, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	if button == 1 {
		return robotgo.LeftButton
	} else if button == 3 {
		return robotgo.RightButton
	}
	return 0
}

func getMousePos(s string) (int, int) {
	xStr := s[subIndex(s, "X:")+2 : subIndex(s, "Y:")]
	yStr := s[subIndex(s, "Y:")+2:]
	x, _ := strconv.Atoi(xStr)
	y, _ := strconv.Atoi(yStr)
	return x, y
}

func subIndex(s string, substr string) int {
	index := -1
	if len(substr) == 0 {
		return index
	}
	for i := range s {
		if i+len(substr) > len(s) {
			break
		}
		if s[i:i+len(substr)] == substr {
			index = i
			break
		}
	}
	return index
}

func main() {
	// test for MousePress and MouseRelease functions
	MousePress("1")
	robotgo.MilliSleep(100)
	MouseRelease("1")

	// test for MouseMove function
	MouseMove("X:100 Y:100")
}
