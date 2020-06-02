// Code for 'Battleship' game server
// Written by Yechan Lee, 2020.04. -
// Game version 0.0.0 (on development)

package main

import (
	"fmt"
	"net"
	"os"
)

var (
	numCurrentUser int
	numMaxUser     int
)

func main() {
	fmt.Println("------------------------------------")
	fmt.Println("'Battleship' game server.")
	fmt.Println("Written by Yechan Lee, 2020.04. -")
	fmt.Println("------------------------------------")

	currentUsers := make(map[net.Conn]net.Addr)
	readyUsers := make(map[net.Conn]net.Addr)
	numCurrentUser = 0
	go userCommand(currentUsers)

	listener, err := net.Listen("tcp", ":8200")
	// error handling
	if err != nil {
		fmt.Println("Failed to tcp Listen(): ", err)
		return
	}
	fmt.Println("Network protocol: \"tcp\", Address: \"121.159.177.222\"\nListening: \":8200\"")
	fmt.Print("\n")

	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to tcp Accept(): ", err)
			continue
		} else {
			fmt.Print(connection.RemoteAddr())
			fmt.Print(" joined the server.\n")
			currentUsers[connection] = connection.RemoteAddr()
			numCurrentUser += 1
		}
		defer closeConnection(connection, currentUsers, connection.RemoteAddr())

		go requestHandler(connection, currentUsers, readyUsers)
	}
}

func closeConnection(connection net.Conn, currentUsers map[net.Conn]net.Addr, address net.Addr) { //* not working properly
	var foundUser bool = false
	for conn, addr := range currentUsers { // find address index in 'currentUsers' map
		if addr == address {
			fmt.Println(addr, "left the server.")
			delete(currentUsers, conn)
			foundUser = true

			numCurrentUser -= 1
		}
	}
	if !foundUser { // error handling
		fmt.Println("Error in 'closeConnection': failed to find and erase user address in 'currentUsers' map.")
		return
	}

	connection.Close()
}

func requestHandler(connection net.Conn, currentUsers map[net.Conn]net.Addr, readyUsers map[net.Conn]net.Addr) {
	// received data is stored in it
	data := make([]byte, 4096)

	for {
		// read message
		_, err := connection.Read(data)
		if err != nil {
			fmt.Println("Failed to read data: ", err)
			return
		}

		// reply message
		for conn := range currentUsers {
			_, err = conn.Write(data)
			fmt.Println(connection, "Passed the message to ", conn)
			if err != nil {
				fmt.Println("Failed to write data: ", err)
				return
			}
		}
	}
}

func userCommand(currentUsers map[net.Conn]net.Addr) {
	for {
		var cmd string
		fmt.Scan(&cmd)
		if cmd == "/help" {
			fmt.Println("----------------------------------------")
			fmt.Println("/exit: quit")
			fmt.Println("/user: print current user information")
			fmt.Println("----------------------------------------")
		} else if cmd == "/exit" {
			fmt.Println("server shut down.")
			os.Exit(1)
		} else if cmd == "/user" {
			fmt.Println("----------------------------------------")
			fmt.Print("Number of user(s): ", numCurrentUser, "\n")
			fmt.Println("Network Connection\tAddress")
			for conn, addr := range currentUsers {
				fmt.Println(conn, addr)
			}
			fmt.Println("----------------------------------------")
		} else {
			fmt.Println("Failed to find command: ", cmd)
		}
	}
}
