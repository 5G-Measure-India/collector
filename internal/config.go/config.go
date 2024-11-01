package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/5G-Measure-India/collector/internal/util"
	"github.com/spf13/pflag"
)

var (
	OutDir string
	Python string
	PhyPy  string

	help    bool
	version bool

	Timestamp time.Time

	NAME    string
	VERSION string
)

func Define() {
	pflag.StringVarP(&OutDir, "out-dir", "o", "data", "Output directory")
	pflag.StringVarP(&PhyPy, "python", "p", "python3", "Python command")
	pflag.StringVarP(&PhyPy, "phy-py", "y", "mi.py", "Phy logger python file")

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

	Timestamp = time.Now()
	log.Printf("[config] timestamp: %.9f", util.GetTime(Timestamp))
}
