// Code for 'Battleship' game server
// Written by Yechan Lee, 2020.04. -
// Game version 0.0.0 (on development)

package main

import (
	"fmt"
	"net"
	"os"
)

var ()

func main() {
	currentUsers := make([]net.Addr, 5)
	go userCommand(currentUsers)

	listener, err := net.Listen("tcp", ":8200")

	// error handling
	if err != nil {
		fmt.Print(err)
		return
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Print(err)
			continue
		} else {
			fmt.Print(connection.LocalAddr())
			fmt.Print(" connected to server.")
			currentUsers = append(currentUsers, connection.LocalAddr())
		}
		defer connection.Close()

		go requestHandler(connection)
	}
}

func requestHandler(connection net.Conn) {
	data := make([]byte, 4096)

	for {
		n, err := connection.Read(data)

		if err != nil {
			fmt.Println(err)
			return
		}

		connection.Write(data[:n])
		fmt.Println("Server read the message: " + string(data[:n]))

		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func userCommand(currentUsers []net.Addr) {
	for {
		var cmd string
		fmt.Scan(&cmd)
		if cmd == "/help" {
			fmt.Println("/q: quit")
			fmt.Println("/u: print current user information")
		}
		if cmd == "/q" {
			os.Exit(1)
		} else if cmd == "/user" {
			for _, user := range currentUsers {
				fmt.Println(user)
			}
		}
	}
}
