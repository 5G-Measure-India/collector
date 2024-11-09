package net

import (
	"log"
	"os"
	"path"

	"github.com/5G-Measure-India/collector/internal/config"
	"github.com/5G-Measure-India/collector/internal/util"
)

func PingRoutine() {
	defer close(util.PingDone)

	csvFile := path.Join(config.OutDir, config.GetFname("ping.csv"))
	csvWriter, err := os.Create(csvFile)
	if err != nil {
		log.Println("[ping] error opening log file:", err)
		return
	}
	defer csvWriter.Close()

	if _, err := csvWriter.WriteString("timestamp,rtt\n"); err != nil {
		log.Println("[ping] error writing csv header:", err)
		return
	}

	cmd := config.Adb("shell", path.Join(config.TOOLS_DIR, "tping"), "-i", "100", "-f", "csv", config.PingServer)
	cmd.Stdout = csvWriter

	if err := cmd.Start(); err != nil {
		log.Println("[ping] error starting:", err)
		return
	}
	log.Println("[ping] started | logging to:", csvFile)

	<-util.Stop

	if err := cmd.Process.Kill(); err != nil {
		log.Println("[ping] error stopping:", err)
	} else if _, err := cmd.Process.Wait(); err != nil {
		log.Println("[ping] error stopping:", err)
	}
}
