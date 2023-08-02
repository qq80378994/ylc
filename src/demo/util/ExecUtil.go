package util

import (
	"fmt"
	"os"
	"os/exec"
)

func RestartProgram() {
	// 获取当前可执行文件的路径
	executable, err := os.Executable()
	if err != nil {
		fmt.Println("无法获取可执行文件路径：", err)
		return
	}

	// 构建一个新的命令来重启当前程序
	cmd := exec.Command(executable, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 启动新的进程来替换当前进程
	err = cmd.Start()
	if err != nil {
		fmt.Println("无法启动新进程：", err)
		return
	}

	// 退出当前进程
	os.Exit(0)
}
