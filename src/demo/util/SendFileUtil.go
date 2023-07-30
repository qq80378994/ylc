package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
	"ylc/src/demo/MyConst"
)

var FileWriter io.Writer

var file *os.File // 声明一个全局变量来存储文件对象

// CloseFile 关闭文件
func CloseFile(socket net.Conn, context string) {
	file.Close()
	compress, _ := Compress([]byte(context))
	Send(MyConst.FILE_UPLOAD_END, compress, socket)
}

// CreateFile 创建文件
func CreateFile(filename string) error {
	filename = strings.Replace(filename, "\\", "", 1)
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	file = f // 将创建的文件对象赋值给全局变量
	return nil
}

func PrepareFile() error {
	file, err := os.Create(file.Name())
	if err != nil {
		return err
	}
	FileWriter = file
	return nil
}

func FileUpload(data string) error {
	bytes := []byte(data)

	_, err := FileWriter.Write(bytes)
	if err != nil {
		fmt.Println("1111111111111111111111111111111111111111111111111111")
		fmt.Println(err)
		return err
	}
	return nil
}

func FileDownload(path string, socket net.Conn) {
	path = strings.Replace(path, "\\", "", 1)
	fmt.Println("开始下载")
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	defer file.Close()

	filename := file.Name()
	filenameReal := filepath.Base(filename)
	fmt.Println(filenameReal)

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	compressFileName, err := Compress([]byte(filenameReal))
	err = Send(MyConst.FILE_CREATEWITHNAME, compressFileName, socket)

	fileSize := fileInfo.Size()
	bufferSize := 8 * 1024
	bytes := make([]byte, bufferSize)

	if fileSize < int64(bufferSize) {
		file.Read(bytes[:fileSize])

		SendHead(byte(MyConst.FILE_PREPARE), socket)
		compressFileSize, err := Compress(bytes[:fileSize])
		err = Send(MyConst.FILE_DOWNLOAD, compressFileSize, socket)

		if err != nil {
			fmt.Println(err)
			fmt.Println(err)
		}
		compressFileName, err := Compress([]byte(filename))
		err = Send(MyConst.FILE_DOWNLOAD_END, compressFileName, socket)

	} else {
		SendHead(byte(MyConst.FILE_PREPARE), socket)

		numBuffers := int(fileSize / int64(bufferSize))
		remainingBytes := int(fileSize % int64(bufferSize))

		for i := 0; i < numBuffers; i++ {
			file.Read(bytes)
			compressByte, err := Compress([]byte(bytes))
			err = Send(MyConst.FILE_DOWNLOAD, compressByte, socket)
			if err != nil {
				fmt.Println(err)
				fmt.Println(err)
			}
		}

		if remainingBytes > 0 {
			bytes = make([]byte, remainingBytes)
			file.Read(bytes)
			compressByte, err := Compress([]byte(bytes))
			err = Send(MyConst.FILE_DOWNLOAD, compressByte, socket)
			if err != nil {

				fmt.Println(err)
			}

		}

		compressFileName, err := Compress([]byte(filename))
		fmt.Println("下载结束")
		err = Send(MyConst.FILE_DOWNLOAD_END, compressFileName, socket)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// OpenFile 打开文件
func OpenFile(filename string) error {
	filename = strings.Replace(filename, "\\", "", 1)
	cmd := exec.Command("cmd", "/c", "start", filename)

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}

	return nil
}

// SendFileWindow 显示界面
func SendFileWindow(conn net.Conn) {

	err := SendHead(byte(MyConst.SHOW_FILEWINDOW), conn)
	if err != nil {
		return
	}

}
func GetLogicalDrives() []string {
	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	GetLogicalDrives := kernel32.MustFindProc("GetLogicalDrives")
	n, _, _ := GetLogicalDrives.Call()
	s := strconv.FormatInt(int64(n), 2)
	var drives_all = []string{"A:", "B:", "C:", "D:", "E:", "F:", "G:", "H:", "I:", "J:", "K:", "L:", "M:", "N:", "O:", "P：", "Q：", "R：", "S：", "T：", "U：", "V：", "W：", "X：", "Y：", "Z："}
	temp := drives_all[0:len(s)]
	var d []string
	for i, v := range s {

		if v == 49 {
			l := len(s) - i - 1
			d = append(d, temp[l])
		}
	}
	var drives []string
	for _, v := range d {
		v = v + "\\"
		drives = append(drives, v)
	}
	return drives
}

// DiskQuery 读取根目录
func DiskQuery(socket net.Conn) {
	roots := GetLogicalDrives()
	//roots := filepath.SplitList(os.Getenv("PATH"))
	for _, root := range roots {
		disk := filepath.VolumeName(root)
		info := fmt.Sprintf("Disk|%s| | ", disk+"\\")
		compress, err := Compress([]byte(info))
		err = Send(MyConst.FILE_QUERY, compress, socket)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// FileQuery 读取信息
func FileQuery(socket net.Conn, path string) {
	// 使用 strings.Replace() 函数将第一个 "\\" 替换为空字符串
	path = strings.Replace(path, "\\", "", 1)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		filePath := path + "/" + f.Name()
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Println("Failed to get file info for", filePath, "-", err)
			continue
		}

		var info string
		if fileInfo.IsDir() {
			fmt.Println("=====131")
			info = fmt.Sprintf("Directory|%s|%s| ", f.Name(), fileInfo.ModTime().Format(time.RFC3339))
		} else {
			fmt.Println("=====132")
			info = fmt.Sprintf("File|%s|%s|%s", f.Name(), fileInfo.ModTime().Format(time.RFC3339), getFileSizeCompany(fileInfo.Size()))
		}
		compress, err := Compress([]byte(info))
		err = Send(MyConst.FILE_QUERY, compress, socket)
	}

}

// 获取文件大小
func getFileSizeCompany(size int64) string {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.2fTB", float64(size)/TB)
	case size >= GB:
		return fmt.Sprintf("%.2fGB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2fMB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2fKB", float64(size)/KB)
	default:
		return fmt.Sprintf("%dB", size)
	}
}
