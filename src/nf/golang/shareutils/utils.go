package shareutils

import (
	"os"
	"flag"
	"fmt"
	"log"
)

type ErrMessage struct {
	Error string
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