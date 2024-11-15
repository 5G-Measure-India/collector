package config

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/5G-Measure-India/collector/internal/util"
	"github.com/spf13/pflag"
)

var (
	OutDir     string
	adb        string
	Python     string
	PingServer string

	help    bool
	version bool

	timestamp time.Time
)

const (
	NAME      = "collector"
	VERSION   = "v0.1.0"
	TOOLS_DIR = "/data/local/tmp"
	DATA_DIR  = "/sdcard/collector"
)

func Define() {
	pflag.StringVarP(&OutDir, "out-dir", "o", "data", "Output directory")
	pflag.StringVarP(&adb, "adb", "a", "adb", "Path to adb")
	pflag.StringVarP(&Python, "python", "p", "python3", "Path to python")
	pflag.StringVarP(&PingServer, "ping-server", "s", "1.1.1.1", "Ping server")

	pflag.BoolVarP(&help, "help", "h", false, "Show this help")
	pflag.BoolVarP(&version, "version", "v", false, "Show version")

	pflag.CommandLine.SortFlags = false
}

func Parse() {
	pflag.Parse()

	if help {
		fmt.Printf("\nUsage: %s [OPTIONS]\n\nOptions:\n", NAME)
		pflag.PrintDefaults()
		os.Exit(0)
	}

	if version {
		fmt.Printf("%s %s\n", NAME, VERSION)
		os.Exit(0)
	}

	verifyFlags()

	timestamp = time.Now()
	log.Printf("[config] timestamp: %.9f", util.GetTime(timestamp))
}

func GetFname(fileName string) string {
	ext := path.Ext(fileName)

	return fileName[:len(fileName)-len(ext)] + "-" + timestamp.Format("2006-01-02T15-04-05Z0700") + ext
}

func Adb(arg ...string) *exec.Cmd {
	return exec.Command(adb, arg...)
}
