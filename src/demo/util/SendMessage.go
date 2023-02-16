package util

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"net"
)

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
