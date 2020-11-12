package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/sunshineplan/utils/httpsvr"
	"github.com/vharitonsky/iniflags"
)

var self string
var logPath string
var server httpsvr.Server

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
	flag.StringVar(&server.Unix, "unix", "", "UNIX-domain Socket")
	flag.StringVar(&server.Host, "host", "0.0.0.0", "Server Host")
	flag.StringVar(&server.Port, "port", "12345", "Server Port")
	//flag.StringVar(&logPath, "log", "/var/log/app/sda-go.log", "Log Path")
	flag.StringVar(&logPath, "log", "", "Log Path")
	iniflags.SetConfigFile(filepath.Join(filepath.Dir(self), "config.ini"))
	iniflags.SetAllowMissingConfigFile(true)
	iniflags.Parse()

	run()
}
