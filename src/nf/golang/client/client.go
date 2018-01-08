//package main
//
//import (
//	"bufio"
//	"fmt"
//	"net"
//)
//
//func main() {
//	p := make([]byte, 2048)
//	conn, err := net.Dial("udp", "127.0.0.1:1234")
//	if err != nil {
//		fmt.Printf("Some error %v", err)
//		return
//	}
//	fmt.Fprintf(conn, "Hi UDP Server, How are you doing?")
//	_, err = bufio.NewReader(conn).Read(p)
//	if err == nil {
//		fmt.Printf("%s\n", p)
//	} else {
//		fmt.Printf("Some error %v\n", err)
//	}
//	conn.Close()
//}

/*
Implements the solution to assignment 1 for UBC CS 416 2017 W2.

Usage:
$ go run client.go [local UDP ip:port] [local TCP ip:port] [aserver UDP ip:port]

Example:
$ go run client.go 127.0.0.1:2020 127.0.0.1:3030 127.0.0.1:7070

*/

package main

import (
	"crypto/md5"
	"encoding/hex"
	//"nf/golang/netWorkutils"
	"nf/golang/netWorkutils"
)

/////////// Msgs used by both auth and fortune servers:

// An error message from the server.
type ErrMessage struct {
	Error string
}

/////////// Auth server msgs:

// Message containing a nonce from auth-server.
type NonceMessage struct {
	Nonce string
	N     int64 // PoW difficulty: number of zeroes expected at end of md5(nonce+secret)
}

// Message containing an the secret value from client to auth-server.
type SecretMessage struct {
	Secret string
}

// Message with details for contacting the fortune-server.
type FortuneInfoMessage struct {
	FortuneServer string // TCP ip:port for contacting the fserver
	FortuneNonce  int64
}

/////////// Fortune server msgs:

// Message requesting a fortune from the fortune-server.
type FortuneReqMessage struct {
	FortuneNonce int64
}

// Response from the fortune-server containing the fortune.
type FortuneMessage struct {
	Fortune string
	Rank    int64 // Rank of this client solution
}

// Main workhorse method.
func main() {
	// TODO
	// Use json.Marshal json.Unmarshal for encoding/decoding to servers

	//UDP test
	udp := netWorkutils.UdpConnection{"198.162.33.54:5555","hello"}
	netWorkutils.Connect(udp)

	////TCP test
	//tcp := netWorkutils.TcpConnection{"127.0.0.1:1234","bye"}
	//netWorkutils.Connect(tcp)

	//fmt.Print(shareutils.ParseCommandLine())



}

// Returns the MD5 hash as a hex netWorkutils for the (nonce + secret) value.
func computeNonceSecretHash(nonce string, secret string) string {
	h := md5.New()
	h.Write([]byte(nonce + secret))
	str := hex.EncodeToString(h.Sum(nil))
	return str
}

