package trans

import (
	"log"
	"os"
	"os/exec"

	"github.com/5G-Measure-India/collector/internal/channel"
	"github.com/5G-Measure-India/collector/internal/config.go"
	"github.com/5G-Measure-India/collector/internal/util"
)

func PingRoutine() {
	defer close(channel.PingDone)

	csvFile := util.GetFilePath(config.OutDir, "ping.csv", config.Timestamp)
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

	cmd := exec.Command("adb", "shell", "/data/local/tmp/tping", "-S", "8.8.8.8", "-s", "1000", "-f", "csv")
	cmd.Stdout = csvWriter

	if err := cmd.Start(); err != nil {
		log.Println("[ping] error starting:", err)
		return
	}
	log.Println("[ping] started | logging to:", csvFile)

	<-channel.Stop

	if err := cmd.Process.Kill(); err != nil {
		log.Println("[ping] error stopping:", err)
	} else if _, err := cmd.Process.Wait(); err != nil {
		log.Println("[ping] error stopping:", err)
	}
}
