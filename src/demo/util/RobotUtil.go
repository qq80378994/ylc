package util

import (
	"fmt"
	"strconv"

	"github.com/go-vgo/robotgo"
)

func KeyReleased(s string) {
	fmt.Println("释放===》", s)
	robotgo.KeyToggle(s, "up")
}

func KeyPress(s string) {

	uMap := robotgo.Keycode
	keyname := getKeyByValue(s, uMap)
	fmt.Println("按下===》", keycode)
	robotgo.KeyToggle(keycode, "down")
}

func MouseWheel(s string) {
	i, _ := strconv.Atoi(s)
	fmt.Println("滑轮===》", i)
	robotgo.Scroll(0, i)
}
func MousePress(s string) {
	fmt.Println("鼠标点击了===》", s)

	button := getMouseClick(s)
	robotgo.MouseToggle("down", button)
	//robotgo.Click(button, true)

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
	//robotgo.Click(button, false)
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
func getKeyByValue(m map[string]int, v int) string {
	for k, val := range m {
		if val == v {
			return k
		}
	}
	// 如果没有找到匹配的键，则返回空字符串
	return ""
}
func main() {
	// test for MousePress and MouseRelease functions
	MousePress("1")
	robotgo.MilliSleep(100)
	MouseRelease("1")

	// test for MouseMove function
	MouseMove("X:100 Y:100")
}
