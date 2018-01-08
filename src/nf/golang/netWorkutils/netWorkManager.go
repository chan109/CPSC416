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
func (udp UdpConnection) connect(){
	p :=  make([]byte, 2048)
	conn, err := net.Dial("udp", udp.Host)
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
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


