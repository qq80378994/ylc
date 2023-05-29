package util

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"strconv"
	"syscall"
)

// 导入Windows API
var (
	user32DLL       = syscall.NewLazyDLL("user32.dll")
	getKeyNameTextW = user32DLL.NewProc("GetKeyNameTextW")
)

func KeyReleased(keycode string) {
	atoi, _ := strconv.Atoi(keycode)
	char := fmt.Sprintf("%c", atoi)
	fmt.Println("按下===》", char)
	robotgo.KeyToggle(char, "down")
}

func KeyPress(keycode string) {
	atoi, _ := strconv.Atoi(keycode)
	char := fmt.Sprintf("%c", atoi)
	fmt.Println("按下===》", char)
	robotgo.KeyToggle(char, "down")
}

func MouseWheel(s string) {
	i, _ := strconv.Atoi(s)
	fmt.Println("滑轮===》", i)
	robotgo.Scroll(0, i)
}
func MousePress(s string) {
	button := getMouseClick(s)
	robotgo.MouseToggle("down", button)

}

func MouseDragged(s string) {
	fmt.Println("鼠标拖拽了===》", s)
	button := getMouseClick(s)
	robotgo.MouseToggle("down", button)
}

func MouseRelease(s string) {
	fmt.Println("鼠标释放了")
	button := getMouseClick(s)
	robotgo.MouseToggle("up", button)
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
	} else if button == 3 {
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
func getKeyByValue(uMap map[string]uint16, value uint16) (string, bool) {
	for key, val := range uMap {
		if val == value {
			return key, true
		}
	}
	return "", false
}
func main() {
	// test for MousePress and MouseRelease functions
	KeyPress("65")
}
