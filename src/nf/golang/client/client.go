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
	"nf/golang/netWorkutils"
	"nf/golang/shareutils"
	"encoding/json"
	"log"
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

func throwError(){
	panic("THere is an error")
}

// Main workhorse method.
func main() {
	var args[]string = shareutils.ParseCommandLine()

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
	udp := netWorkutils.UdpConnection{host,msg, localAddr}
	return udp.Connect()
}

//send message through tcp to server
func sendMsgTcp(host string, msg string, localAddr string) string  {
	tcp := netWorkutils.TcpConnection{host,msg, localAddr}
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

