//package main
//
//import "fmt"
//
//func allPossibleCombinations(input string, length int, curstr string) []string {
//	if(len(curstr) == length) {
//		return[]string{curstr};
//	}
//
//	var ret []string
//	for i := 0; i < len(input); i++ {
//		allPossibleCombinations(input, length, curstr + string(input[i]));
//	}
//	return ret;
//}
//
//
//
//func main() {
//	//input :=[]string{"a", "b", "c", "d"}
//	input :="abcd"
//
//	fmt.Println(allPossibleCombinations(input, 3, ""));
//}

/*
Implements the solution to assignment 1 for UBC CS 416 2016 W2.

Usage:
$ go run client.go [local UDP ip:port] [server UDP ip:port]

Example:
$ go run client.go 127.0.0.1:2020 127.0.0.1:7070

*/

package main

import (
	"encoding/gob"
	"fmt"
	"os"

	//Stewarts Imports
	"bytes"
	"log"
	"net"
	"time"
)

//local_ip_port := "128.189.118.124:8888"
//
//l2 := "198.162.33.54:5555"

// Main workhorse method.
func main() {

	// Missing command line args.

	// Extract the command line args.
	local_ip_port := "128.189.118.124:8888"
	remote_ip_port := "142.103.15.6:7777"

	//remote_ip_port := "198.162.33.54:5555"

	//Stewarts Solution
	logger = log.New(os.Stdout, "[416-A1-Stewarts-Solution] ", log.Lshortfile)
	conn := getConnection(local_ip_port)
	server := getAddr(remote_ip_port)

	min := uint32(0)
	max := uint32(0xFFFFFFFF)
	buffer := make([]byte, 1024)

	for min <= max {
		mid := min/2 + max/2
		payload, err := Marshall(mid)
		if err != nil {
			logger.Fatal(err)
		}
		conn.WriteToUDP(payload, server)
		conn.SetDeadline(time.Now().Add(time.Millisecond))
		n, err := conn.Read(buffer)
		fmt.Println(string(buffer[:bytes.Index(buffer, []byte{0})]))
		fmt.Println(string(n))
		if err != nil {
			continue
		}
		response := string(buffer[0:n])
		if response == "low" {
			min = mid
		} else if response == "high" {
			max = mid
		} else {
			logger.Println(response)
			os.Exit(0)
		}
	}
	logger.Println("Unable to find the key to fortune; time to die")
}

func Marshall(guess uint32) ([]byte, error) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(guess)
	return network.Bytes(), err
}

//listen connection returns a udp conn for listening on
func getConnection(ip string) *net.UDPConn {
	lAddr, err := net.ResolveUDPAddr("udp", ip)
	if err != nil {
		logger.Fatal(err)
	}
	l, err := net.ListenUDP("udp", lAddr)
	if err != nil {
		logger.Fatal(err)
	}
	return l
}

func getAddr(ip string) *net.UDPAddr {
	addr, err := net.ResolveUDPAddr("udp", ip)
	if err != nil {
		logger.Fatal(err)
	}
	return addr
}

var (
	logger *log.Logger
)