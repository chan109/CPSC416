package netWorkutils

import (
	"fmt"
	"net"
	"bufio"
	"bytes"
)

type connection interface {
	connect()
}

// Export udp struct definition
type UdpConnection struct{
	Host string
	Msg string
}

// Export tcp struct definition
type TcpConnection struct{
	Host string
	Msg string
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

//private method for udp connection(accept local and remote address)
func (udp UdpConnection) connect(){
	sip, err := net.ResolveUDPAddr("udp","198.162.33.23:8888")
	CheckError(err)
	//checkerror
	dip, err := net.ResolveUDPAddr("udp",udp.Host)
	CheckError(err)
	//checkerror

	conn, err := net.DialUDP("udp", sip, dip)
	p :=  make([]byte, 2048)

	fmt.Fprintf(conn, udp.Msg)
	_, err = bufio.NewReader(conn).Read(p)
	if err == nil {
		fmt.Printf("%s\n", p[:bytes.Index(p, []byte{0})])
	} else {
		fmt.Printf("Some error %v\n", err)
	}
	conn.Close()


}

//private method for tcp connection
func (tcp TcpConnection) connect(){

	p := make([]byte, 2048)

	// connect to this socket
	conn, err := net.Dial("tcp", tcp.Host)
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	// send to socket
	fmt.Fprintf(conn, tcp.Msg + "\n")
	//losten for reply
	_, err = bufio.NewReader(conn).Read(p)
	if err == nil {
		fmt.Printf("%s\n", p[:bytes.Index(p, []byte{0})])
	} else {
		fmt.Printf("Some error %v\n", err)
	}
	conn.Close()

}

// Export method for TCP/UP connection
func Connect(obj connection)  {
	obj.connect()
}


