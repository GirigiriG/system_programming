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
		fmt.Println("ECONNREFUSED to localhost:"+port, err)
		os.Exit(111)
	}
	defer conn.Close()

	go func() {
		// buf := make([]byte, 1024)
		reader := bufio.NewReader(conn)
		for {
			msg, err := reader.ReadString('\n') // or some delimiter

			if err != nil {
				fmt.Println("connection closed by server")
				os.Exit(1)
				return // don't exit, just stop goroutine
			}

			fmt.Print(msg)
		}
	}()

	stdinScanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("enter message: ")
		if stdinScanner.Scan() {
			msg := stdinScanner.Text()
			if _, err := conn.Write([]byte(msg+"\n")); err != nil {
				fmt.Println("Error writing to server:", err)
			}
		}
	}
}
