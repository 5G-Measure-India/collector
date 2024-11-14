package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/5G-Measure-India/collector/internal/config"
	"github.com/5G-Measure-India/collector/internal/net"
	"github.com/5G-Measure-India/collector/internal/phy"
	"github.com/5G-Measure-India/collector/internal/util"
)

func main() {
	config.Define()
	config.Parse()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go phy.GeoRoutine()
	go phy.PhyRoutine()
	go net.PingRoutine()
	go net.SpeedtestRoutine()

	<-done
	println()

	close(util.Stop)

	<-util.SpeedtestDone
	<-util.PingDone
	<-util.PhyDone
	<-util.GeoDone
}
