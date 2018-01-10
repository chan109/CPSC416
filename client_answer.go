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
	"math/rand"
	"fmt"
	"time"
	"encoding/json"
	"log"
	"bufio"
	"bytes"
	"net"
	"flag"
	"os"
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

type UdpConnection struct{
	Host string
	Msg string
	LocalAddr string
}

// Export tcp struct definition
type TcpConnection struct{
	Host string
	Msg string
	LocalAddr string
}

func CheckError(err error) {
	if err  != nil {
		fmt.Println("Error: " , err)
	}
}

func readData(conn *net.UDPConn, err error) string{
	p :=  make([]byte, 2048)

	_, err = bufio.NewReader(conn).Read(p)
	if err == nil {
		conn.Close()
		return string(p[:bytes.Index(p, []byte{0})])
	} else {
		fmt.Printf("Some error %v\n", err)
		conn.Close()
		return "-1"
	}
}

//private method for udp connection(accept local and remote address)
//udp local address has to be the public one "128.189.112.244:8888" or run curl ipinfo.io/ip to help getting the public ip
func (udp UdpConnection) Connect() string{

	p :=  make([]byte, 2048)

	//get local ip
	sip, err := net.ResolveUDPAddr("udp",udp.LocalAddr)
	CheckError(err)

	//get remote ip
	dip, err := net.ResolveUDPAddr("udp",udp.Host)
	CheckError(err)

	conn, err := net.DialUDP("udp", sip, dip)
	CheckError(err)

	fmt.Fprintf(conn, udp.Msg)
	_, err = bufio.NewReader(conn).Read(p)
	if err == nil {
		conn.Close()
		fmt.Println(string(p[:bytes.Index(p, []byte{0})]))
		return string(p[:bytes.Index(p, []byte{0})])

	} else {
		fmt.Printf("Some error %v\n", err)
		conn.Close()
		return "-1"
	}
}

//private method for tcp connection
func (tcp TcpConnection) Connect() string{

	sip, err := net.ResolveTCPAddr("tcp",tcp.LocalAddr)
	CheckError(err)
	//checkerror
	dip, err := net.ResolveTCPAddr("tcp",tcp.Host)
	CheckError(err)

	// connect to this socket
	conn, err := net.DialTCP("tcp", sip, dip)
	CheckError(err)
	// send to socket

	p :=  make([]byte, 2048)

	fmt.Fprintf(conn, tcp.Msg)
	_, err = bufio.NewReader(conn).Read(p)
	if err == nil {
		//fmt.Printf("%s\n", p[:bytes.Index(p, []byte{0})])
		conn.Close()
		return string(p[:bytes.Index(p, []byte{0})])

	} else {
		fmt.Printf("Some error %v\n", err)
		conn.Close()
		return "-1"
	}

}


func throwError(){
	panic("THere is an error")
}

func ParseCommandLine()[]string  {
	flag.Usage = func() {
		fmt.Println("Usage of the program:")
		fmt.Printf("go run client.go [local UDP ip:port] [local TCP ip:port] [aserver UDP ip:port]\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := os.Args
	if(len(args)!=4){
		log.Fatal("Argumens has to be 3")
		return []string{}
	}

	return []string{args[1], args[2], args[3]}

}
// Main workhorse method.
func main() {
	var args[]string = ParseCommandLine()

	//get the nonece msg through UDP
	var str = sendMsgUdp(args[2],"hello", args[0])
	var nonceMsg NonceMessage
	var fserverIp string
	var fortuneNonce int64

	//parse the nonce
	if(str == "-1"){
		throwError()
	}else{
		if err := json.Unmarshal([]byte(str), &nonceMsg); err != nil{
			log.Fatal(err)
		}
	}

	//find the secret using the parsed nonce
	secret :=findSecret(nonceMsg.Nonce,nonceMsg.N)

	encodedSecret:= SecretMessage{secret}
	bUdp, err := json.Marshal(encodedSecret)
	if err != nil{
		throwError()
	}

	var fortuneInfoMsg = sendMsgUdp(args[2],string(bUdp), args[1])
	var decodedFortuneMsg FortuneInfoMessage
	json.Unmarshal([]byte(fortuneInfoMsg), &decodedFortuneMsg)
	fserverIp = decodedFortuneMsg.FortuneServer
	fortuneNonce = decodedFortuneMsg.FortuneNonce

	//TCP test
	endCodedFortunce := FortuneReqMessage{fortuneNonce}
	bTcp, _:=json.Marshal(endCodedFortunce)

	answer := sendMsgTcp(fserverIp,string(bTcp), args[1])
	var decodedAnswer FortuneMessage
	json.Unmarshal([]byte(answer), &decodedAnswer)
	fmt.Println(decodedAnswer.Fortune)

}

//send message through udp to server
func sendMsgUdp(host string, msg string, localAddr string) string  {
	udp := UdpConnection{host,msg, localAddr}
	return udp.Connect()
}

//send message through tcp to server
func sendMsgTcp(host string, msg string, localAddr string) string  {
	tcp := TcpConnection{host,msg, localAddr}
	return tcp.Connect()
}

//if no solution is found, it runs forever
func findSecret(nonece string, N int64) string{
	for{
		var valToCompute string = RandStringRunes(8)
		var computedHash string = computeNonceSecretHash(nonece, valToCompute)
		if(checkHash(N, computedHash)){
			//fmt.Println("Found the valid hash: %s", string(computedHash))
			//fmt.Println("Found the valid secret: %s", valToCompute)
			return valToCompute
		}
	}
}

//check the N zeros at the end of the computed hash
func checkHash(N int64, hash string) bool{
	for i := int64(len(hash) -1); i>int64(len(hash))-N-1; i--{
		if(string(hash[i]) != "0") {
			return false
		}
	}

	//check if the position N+1 is zeros of not
	if(string(hash[int64(len(hash))-N-1]) != "0"){
		return true
	}else{
		return false
	}
}

//Generate random string
var letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// Returns the MD5 hash as a hex netWorkutils for the (nonce + secret) value.
func computeNonceSecretHash(nonce string, secret string) string {
	h := md5.New()
	h.Write([]byte(nonce + secret))
	str := hex.EncodeToString(h.Sum(nil))
	return str
}

