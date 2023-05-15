package util

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

func Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := zw.Write(data); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Decompression(context []byte) ([]byte, error) {
	byteReader := bytes.NewReader(context)
	gzipReader, err := gzip.NewReader(byteReader)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer gzipReader.Close()

	byteBuffer := bytes.NewBuffer(nil)
	_, err = byteBuffer.ReadFrom(gzipReader)
	if err != nil {
		return nil, err
	}

	return byteBuffer.Bytes(), nil
}
