package netWorkutils

import (
	"fmt"
	"net"
	"bufio"
	"bytes"
	//"encoding/json"
	//"log"
)

//one way to define interface
//type connection interface {
//	Connect()
//}

// Export udp struct definition
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

type FortuneReqMessage struct {
	FortuneNonce int64
}

//private method for udp connection
//func (udp UdpConnection) connect(){
//	p :=  make([]byte, 2048)
//	conn, err := net.Dial("udp", udp.Host)
//	if err != nil {
//		fmt.Printf("Some error %v", err)
//		return
//	}
//	fmt.Fprintf(conn, udp.Msg)
//	_, err = bufio.NewReader(conn).Read(p)
//	if err == nil {
//		fmt.Printf("%s\n", p[:bytes.Index(p, []byte{0})])
//	} else {
//		fmt.Printf("Some error %v\n", err)
//	}
//	conn.Close()
//}

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



