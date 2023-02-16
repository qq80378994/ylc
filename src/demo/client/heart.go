package client

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func heartRun(socket net.Conn) {
	defer socket.Close()
	writer := bufio.NewWriter(socket)
	for {
		time.Sleep(5 * time.Second)
		fmt.Println(writer, 1)
		writer.Flush()
	}
}
