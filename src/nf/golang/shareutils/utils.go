package shareutils

import (
	"os"
)

type ErrMessage struct {
	Error string
}

func ParseCommandLine()[]string  {
	args := os.Args
	if(len(args)!=4){
		panic("Argumens has to be 3")
		return []string{}
	}
	return []string{args[1], args[2], args[3]}

}