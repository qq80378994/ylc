package util

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"net"
)

func ReceiveHead(dataInputStream *bufio.Reader) byte {
	bytes := make([]byte, 1)
	_, err := dataInputStream.Read(bytes)
	if err != nil {
		// 处理异常
	}
	return bytes[0]
}

func readByte(dataInputStream *bytes.Buffer) (byte, error) {
	var b uint8
	err := binary.Read(dataInputStream, binary.BigEndian, &b)
	if err != nil {
		return 0, err
	}
	return b, nil
}

func Send(head byte, context []byte, conn net.Conn) {
	var buffer bytes.Buffer
	writer := zlib.NewWriter(&buffer)
	writer.Write(context)
	writer.Close()
	data := buffer.Bytes()

	length := uint16(len(data))
	header := []byte{head, byte(length >> 8), byte(length)}
	packet := append(header, data...)

	conn.Write(packet)
}

func ToByte(args ...interface{}) []byte {
	var buffer bytes.Buffer
	for _, arg := range args {
		binary.Write(&buffer, binary.LittleEndian, arg)
	}
	return buffer.Bytes()
}
