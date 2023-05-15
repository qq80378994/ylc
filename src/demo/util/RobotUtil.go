package util

import (
	"fmt"
	"strconv"

	"github.com/go-vgo/robotgo"
)

func MousePress(s string) {
	fmt.Println("鼠标点击了===》", s)
	button := getMouseClick(s)
	robotgo.Toggle("down", button)
}

func MouseRelease(s string) {
	button := getMouseClick(s)
	robotgo.Toggle("up", button)
}

func MouseMove(s string) {
	x, y := getMousePos(s)
	robotgo.MoveSmooth(x, y)
}

func getMouseClick(s string) string {
	button, err := strconv.Atoi(s)
	if err != nil {
		return ""
	}
	if button == 1 {

		return "left"
	} else if button == 2 {
		return "right"
	}
	return ""
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
