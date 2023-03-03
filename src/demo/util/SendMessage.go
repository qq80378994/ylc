package util

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"net"
)

func ReceiveHead(dataInputStream *bufio.Reader) byte {
	bytes := make([]byte, 1)
	_, err := dataInputStream.Read(bytes)
	if err != nil {
		fmt.Println(err)
		// 处理异常
	}
	return bytes[0]
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

func Send(head byte, context []byte, conn net.Conn) {
	//compressed := compress(context)
	fmt.Println(len(context))
	bytes, err := ToByte(head, len(context), context)
	if err != nil {
		fmt.Println("Error converting to bytes:", err)
		return
	}

	_, err = conn.Write(bytes)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}
}

func compress(data []byte) []byte {
	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	zw.Write(data)
	zw.Close()
	return buf.Bytes()
}

func SendHead(head byte, socket net.Conn) {
	buffer := []byte{head}
	_, err := socket.Write(buffer)
	if err != nil {
		fmt.Println(err)
	}
}
