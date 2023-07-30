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
		fmt.Println(err)
		return nil, err
	}
	return b, nil
}
func ReceiveLength(dataInputStream *bufio.Reader) int {
	bytes := make([]byte, 4)
	_, err := dataInputStream.Read(bytes)
	if err != nil {

		fmt.Println(err)
		// 处理异常
	}
	length := ByteToInt(bytes)
	return length
}
func ReceiveHead(dataInputStream *bufio.Reader) (byte, error) {
	bytes := make([]byte, 1)
	_, err := dataInputStream.Read(bytes)
	if err != nil {

		fmt.Println(err)
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
func SendPress(head byte, context []byte, conn net.Conn) error {
	compress, err := Compress(context)
	fmt.Println(len(compress))
	bytes, err := ToByte(head, len(compress), context)
	if err != nil {
		fmt.Println("Error converting to bytes:", err)
	}

	_, err = conn.Write(bytes)
	if err != nil {
		fmt.Println("Error sending data:", err)
	}
	return err
}
func Send(head byte, context []byte, conn net.Conn) error {
	//compressed := compress(context)
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
