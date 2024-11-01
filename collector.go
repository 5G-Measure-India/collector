package main

import (
	"os"

	"github.com/5G-Measure-India/collector/internal/channel"
	"github.com/5G-Measure-India/collector/internal/config.go"
	"github.com/5G-Measure-India/collector/internal/phy"
	"github.com/5G-Measure-India/collector/internal/trans"
)

var sigs = make(chan os.Signal)

func main() {
	config.Define()
	config.Parse()

	go phy.PhyRoutine()
	go trans.PingRoutine()
	go trans.SpeedtestRoutine()

	<-sigs

	channel.StopAll()
}
