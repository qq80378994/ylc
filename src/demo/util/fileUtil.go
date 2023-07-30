package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CopyToProgramData() error {
	src, err := os.Executable()
	if err != nil {
		return err
	}

	dst := filepath.Join("C:\\ProgramData", filepath.Base(src))

	// 判断目标目录是否存在该程序
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		// 目标目录不存在该程序，进行复制操作
		if err := copyFile(src, dst); err != nil {
			return err
		}
	} else {
		// 目标目录已经存在该程序，不进行复制操作
		fmt.Println(dst, "already exists, no need to copy")
	}
	return nil
}

func copyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	return err
}
