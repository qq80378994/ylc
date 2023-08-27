package util

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

var HEART = 99

func ReadPacket(conn net.Conn) ([]byte, error) {
	// 先读取 4 字节，该部分包含了整个数据包的长度
	lengthBytes := make([]byte, 4)
	_, err := io.ReadFull(conn, lengthBytes)
	if err != nil {
		return nil, err
	}

	// 解析出数据包的长度
	packetLength := int(binary.BigEndian.Uint32(lengthBytes))

	// 读取剩余的数据包内容
	packetBytes := make([]byte, packetLength)
	_, err = io.ReadFull(conn, packetBytes)
	if err != nil {
		return nil, err
	}

	return packetBytes, nil
}

func ByteToInt(byte []byte) int {
	bytesBuffer := bytes.NewBuffer(byte)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}

func ReceiveContext(dataInputStream *bufio.Reader, len int) ([]byte, error) {
	bytes := make([]byte, len)
	_, err := dataInputStream.Read(bytes)
	// 获取内容
	b, err := Decompression(bytes)
	if err != nil {
		fmt.Println("请求内容异常", err)
		return nil, err
	}
	return b, nil
}
func ReceiveLength(dataInputStream *bufio.Reader) int {
	bytes := make([]byte, 4)
	_, err := dataInputStream.Read(bytes)
	if err != nil {

		fmt.Println("请求长度异常", err)
		// 处理异常
	}
	length := ByteToInt(bytes)
	return length
}
func ReceiveHead(dataInputStream *bufio.Reader) (byte, error) {
	bytes := make([]byte, 1)
	_, err := dataInputStream.Read(bytes)
	if err != nil {
		fmt.Println("请求头异常", err)
		// 处理异常
	}
	return bytes[0], err
}

func ToByte(head byte, length int, context []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := buf.WriteByte(head)
	if err != nil {
		return nil, err
	}

	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, uint32(length))

	_, err = buf.Write(lengthBytes)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(context)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Send(head byte, context []byte, conn net.Conn) error {
	// 计算消息内容的长度
	// 计算消息内容的长度
	var length uint32
	if context != nil {
		length = uint32(len(context))
	}

	//	fmt.Println("length===>", length)
	// 创建一个字节缓冲区
	buf := new(bytes.Buffer)

	// 使用大端字节序将长度写入缓冲区
	binary.Write(buf, binary.BigEndian, length)

	// 将头部、长度字段和内容合并成一个字节数组
	message := append([]byte{head}, buf.Bytes()...)
	message = append(message, context...)

	// 发送消息
	_, err := conn.Write(message)
	if err != nil {
		fmt.Println("Error sending data:", err)
	}

	return err
}
func SendT(head byte, context []byte, conn net.Conn) error {
	fmt.Println(len(context))
	bytes, err := ToByte(head, len(context), context)
	if err != nil {
		fmt.Println("Error converting to bytes:", err)
	}

	_, err = conn.Write(bytes)
	if err != nil {
		fmt.Println("Error sending data:", err)
	}
	return err
}

func SendHead(head byte, socket net.Conn) error {
	buffer := []byte{head}
	_, err := socket.Write(buffer)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
