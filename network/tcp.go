package network

import (
	"bufio"
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

	reader := bufio.NewReader(conn)
	remoteID := ""

	if strings.Contains(remoteAddress, "::") {
		remoteID = strings.Split(remoteAddress, "]:")[1]
	} else {
		remoteID = strings.Split(remoteAddress, ":")[1]
	}

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("connection closed by client:", remoteAddress)
			} else {
				fmt.Println("read error:", err)
			}
			return
		}

		mu.Lock()
		for addr, con := range connections {
			if addr != remoteAddress {
				fmt.Fprintf(con, "[%v] %v", remoteID, msg)
			}
		}
		mu.Unlock()
	}
}
