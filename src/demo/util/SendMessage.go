package util

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"net"
)

func Compress(context []byte) []byte {
	var buffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&buffer)
	gzipWriter.Write(context)
	gzipWriter.Close()
	return buffer.Bytes()
}

func ReceiveHead(dataInputStream *bufio.Reader) byte {
	bytes := make([]byte, 1)
	_, err := dataInputStream.Read(bytes)
	if err != nil {
		fmt.Println(err)
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

//func Send(head byte, context []byte, socket net.Conn) {
//	compressedData, err := compress(context)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	header := createHeader(head, len(compressedData))
//	_, err = socket.Write(header)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	_, err = socket.Write(compressedData)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//}

func createHeader(head byte, length int) []byte {
	header := make([]byte, 3)
	header[0] = head
	binary.BigEndian.PutUint16(header[1:], uint16(length))
	return header
}
func compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	compressor, err := zlib.NewWriterLevel(&buf, zlib.BestCompression)
	if err != nil {
		return nil, err
	}
	_, err = compressor.Write(data)
	if err != nil {
		return nil, err
	}
	err = compressor.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func SendHead(head byte, socket net.Conn) {
	buffer := []byte{head}
	_, err := socket.Write(buffer)
	if err != nil {
		fmt.Println(err)
	}
}
func ToByte(args ...interface{}) []byte {
	var buffer bytes.Buffer
	for _, arg := range args {
		binary.Write(&buffer, binary.LittleEndian, arg)
	}
	return buffer.Bytes()
}
