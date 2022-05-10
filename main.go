package main

import (
	"os"

	"github.com/0xsuk/byodns/app/controller"
	"github.com/0xsuk/byodns/util"
)

func main() {
	go controller.StartDNSServer()
	controller.StartHTTPServer()
}

type Command struct {
	Name    string
	ArgsLen int
	Usage   string
	Run     func(args ...interface{})
}

func init() {
	if len(os.Args) == 1 {
		util.Println("Starting byodns...")
		return
	}
}
