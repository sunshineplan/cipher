package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/vharitonsky/iniflags"
)

// OS is the running program's operating system
const OS = runtime.GOOS

var self string
var unix, host, port, logPath *string

func init() {
	var err error
	self, err = os.Executable()
	if err != nil {
		log.Fatalf("Failed to get self path: %v", err)
	}
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	fmt.Println(`
  run
		run Simple Text Encryption web service mode
  client
		run Simple Text Encryption client mode (Default)`)
}

func main() {
	flag.Usage = usage
	unix = flag.String("unix", "", "UNIX-domain Socket")
	host = flag.String("host", "127.0.0.1", "Server Host")
	port = flag.String("port", "12345", "Server Port")
	logPath = flag.String("log", "", "Log Path")
	iniflags.SetConfigFile(filepath.Join(filepath.Dir(self), "config.ini"))
	iniflags.SetAllowMissingConfigFile(true)
	iniflags.Parse()

	switch flag.NArg() {
	case 0:
		client()
	case 1:
		switch flag.Arg(0) {
		case "run":
			run()
		case "client":
			client()
		default:
			log.Fatalf("Unknown argument: %s", flag.Arg(0))
		}
	default:
		log.Fatalf("Unknown arguments: %s", strings.Join(flag.Args(), " "))
	}
}
