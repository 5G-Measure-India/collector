package main

import (
	"os"

	"github.com/5G-Measure-India/collector/internal/config"
	"github.com/5G-Measure-India/collector/internal/net"
	"github.com/5G-Measure-India/collector/internal/phy"
	"github.com/5G-Measure-India/collector/internal/util"
)

var sigs = make(chan os.Signal)

func main() {
	config.Define()
	config.Parse()

	go phy.GeoRoutine()
	go phy.PhyRoutine()
	go net.PingRoutine()
	go net.SpeedtestRoutine()

	<-sigs

	close(util.Stop)

	<-util.SpeedtestDone
	<-util.PingDone
	<-util.PhyDone
	<-util.GeoDone
}
