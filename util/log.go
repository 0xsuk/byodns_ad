package util

import (
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const (
	NORMAL_LOG = "/etc/byodns/byodns.log"
	ERR_LOG    = "/etc/byodns/err.log"
)

var (
	INFO *log.Logger
	ERR  *log.Logger
)

func Println(v ...interface{}) {
	INFO.Println(v...)
}

func Printf(s string, v ...interface{}) {
	INFO.Printf(s, v...)
}

func Fatalln(v ...interface{}) {
	s := format(v)
	ERR.Fatalln(s...)
}

func format(v ...interface{}) []interface{} {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	file = file[strings.LastIndex(file, "/")+1:]
	s := []interface{}{file + ":" + strconv.Itoa(line)}
	s = append(s, v...)
	return s
}

func init() {
	//ERR.Fatalln automatically log on to os.Stderr
	n, _ := os.OpenFile(NORMAL_LOG, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	e, _ := os.OpenFile(ERR_LOG, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	nout := io.MultiWriter(os.Stdout, n)
	eout := io.MultiWriter(os.Stdout, n, e)
	INFO = log.New(nout, "[INFO]: ", log.Ldate|log.Ltime)
	ERR = log.New(eout, "\033[31m[ERR]: \033[00m", log.Ldate|log.Ltime)
}
