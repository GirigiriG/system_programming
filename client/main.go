package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	port := "9091"
	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println("unable to connect to localhost:"+port, err)
	}

	defer conn.Close()
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("connection closed by server", err)
			os.Exit(0)
			break
		}

		fmt.Println(string(buf[:n]))
	}

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
