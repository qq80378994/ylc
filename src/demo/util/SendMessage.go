package util

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

var HEART = 99

func ReceiveLength(r io.Reader) (int, error) {
	var length uint32
	err := binary.Read(r, binary.BigEndian, &length)
	return int(length), err
}

// 从给定的读取器中读取给定长度的消息，并将其作为字节切片返回。
func ReceiveContext(r io.Reader, length int) ([]byte, error) {
	context := make([]byte, length)
	_, err := io.ReadFull(r, context)
	return context, err
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

func compress(data []byte) []byte {
	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	zw.Write(data)
	zw.Close()
	return buf.Bytes()
}

func SendHead(head byte, socket net.Conn) error {
	buffer := []byte{head}
	_, err := socket.Write(buffer)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
