package main
import (
	"fmt"
	"net"
	"os"
	"bytes"
)

//TCP
const (
	CONN_HOST = "127.0.0.1"
	CONN_PORT = "1234"
	CONN_TYPE = "tcp"
)


func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
	_,err := conn.WriteToUDP([]byte("From server: Hello I got your mesage "), addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	fmt.Print("message receive:")
	//fmt.Print("%s\n", buf[:bytes.Index(buf, []byte{0})])
	fmt.Print(string(buf[:bytes.Index(buf, []byte{0})]))

	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	// Send a response back to person contacting us.
	conn.Write([]byte("Message received."))
	// Close the connection when you're done with it.
	conn.Close()
}


func main() {

	//UDP
	//p := make([]byte, 2048)
	//addr := net.UDPAddr{
	//	Port: 1234,
	//	IP: net.ParseIP("127.0.0.1"),
	//}
	//ser, err := net.ListenUDP("udp", &addr)
	//if err != nil {
	//	fmt.Printf("Some error %v\n", err)
	//	return
	//}
	//for {
	//	_,remoteaddr,err := ser.ReadFromUDP(p)
	//	fmt.Printf("Read a message from %v %s \n", remoteaddr, p[: bytes.Index(p, []byte{0})])
	//	if err !=  nil {
	//		fmt.Printf("Some error  %v", err)
	//		continue
	//	}
	//	go sendResponse(ser, remoteaddr)
	//}

	//TCP
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}