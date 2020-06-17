// Code for 'Battleship' game server
// Written by Yechan Lee, 2020.04. -
// Game version 1.0.0

package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
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

	// get local IP address
	var currentIP string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		// = GET LOCAL IP ADDRESS
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				currentIP = ipnet.IP.String()
			}
		}
	}
	fmt.Println("Network protocol: \"tcp\", Address: " + currentIP + "\nListening: \":8200\"")
	fmt.Print("\n")

	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to tcp Accept(): ", err)
			continue
		} else {
			fmt.Print("[" + (strconv.Itoa(time.Now().Hour()) + ":" + strconv.Itoa(time.Now().Minute()) + ":" + strconv.Itoa(time.Now().Second()) + "'" + strconv.Itoa(time.Now().Nanosecond())) + "] ")
			fmt.Print(connection.RemoteAddr())
			fmt.Print(" joined the server.\n")
			currentUsers[connection] = connection.RemoteAddr()
			numCurrentUser++
		}
		defer closeConnection(connection, currentUsers, connection.RemoteAddr())

		go requestHandler(connection, currentUsers, readyUsers)
	}
}

func closeConnection(connection net.Conn, currentUsers map[net.Conn]net.Addr, address net.Addr) {
	var foundUser bool
	foundUser = false
	for conn, addr := range currentUsers { // find address index in 'currentUsers' map
		if addr == address {
			fmt.Print("[" + (strconv.Itoa(time.Now().Hour()) + ":" + strconv.Itoa(time.Now().Minute()) + ":" + strconv.Itoa(time.Now().Second()) + "'" + strconv.Itoa(time.Now().Nanosecond())) + "] ")
			fmt.Println(addr, "left the server.")
			delete(currentUsers, conn)
			foundUser = true

			numCurrentUser--
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
	data := make([]byte, 256)

	for {
		// read message
		_, err := connection.Read(data)
		if err != nil { // if tcp connection lost
			fmt.Println("Failed to read data: ", err)
			closeConnection(connection, currentUsers, connection.RemoteAddr())
			return
		}

		// reply message
		for conn := range currentUsers {
			_, err = conn.Write(data)
			fmt.Print("[" + (strconv.Itoa(time.Now().Hour()) + ":" + strconv.Itoa(time.Now().Minute()) + ":" + strconv.Itoa(time.Now().Second()) + "'" + strconv.Itoa(time.Now().Nanosecond())) + "] ")
			fmt.Print(currentUsers[connection], " -> ", currentUsers[conn], "\n")
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
		if cmd == "help" {
			fmt.Println("--------------------------------------------------")
			fmt.Println("exit: quit")
			fmt.Println("user: print current user information")
			fmt.Println("refresh: delete all the data and restart")
			fmt.Println("time: get current time")
			fmt.Println("--------------------------------------------------")
		} else if cmd == "exit" {
			fmt.Println("server shut down.")
			os.Exit(1)
		} else if cmd == "user" {
			fmt.Println("--------------------------------------------------")
			fmt.Print("Number of user(s): ", numCurrentUser, "\n")
			fmt.Println("Index\tNetwork Connection\tAddress")
			var idx int = 0
			for conn, addr := range currentUsers {
				fmt.Println(idx, "\t", conn, "\t", addr)
				idx++
			}
			if idx == 0 {
				fmt.Println("NULL\tNULL\t\t\tNULL")
			}
			fmt.Println("--------------------------------------------------")
		} else if cmd == "refresh" {
			for conn, addr := range currentUsers {
				closeConnection(conn, currentUsers, addr)
			}
		} else if cmd == "time" {
			fmt.Println("Current time: [" + (strconv.Itoa(time.Now().Hour()) + ":" + strconv.Itoa(time.Now().Minute()) + ":" + strconv.Itoa(time.Now().Second()) + "'" + strconv.Itoa(time.Now().Nanosecond())) + "] ")
		} else {
			fmt.Println("Failed to find command: ", cmd)
		}
	}
}
