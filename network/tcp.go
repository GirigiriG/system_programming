package network

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

var (
	connections map[string]net.Conn
	mu          sync.Mutex
)

func init() {
	connections = make(map[string]net.Conn)
}

func New() {
	port := "9091"
	listen, err := net.Listen("tcp", ":"+port)
	fmt.Println("tcp socket opened on " + port)
	if err != nil {
		panic(err.Error())
	}

	for {
		conn, err := listen.Accept()
		fmt.Println("accepted\n" + conn.RemoteAddr().String())
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	remoteAddress := conn.RemoteAddr().String()
	mu.Lock()
	connections[remoteAddress] = conn
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(connections, remoteAddress)
		mu.Unlock()
		conn.Close()
	}()

	buf := make([]byte, 6)

	for {
		_, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("connection closed by client:", remoteAddress)
			} else {
				fmt.Println("read error:", err)
			}
			return
		}

		msg := string(buf)

		mu.Lock()
		for addr, con := range connections {
			if addr != remoteAddress {
				fmt.Fprintf(con, "[%v] %v", strings.Split(addr, "]:")[1], msg)
			}
		}
		mu.Unlock()
	}
}
