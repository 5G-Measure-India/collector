package trans

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/5G-Measure-India/collector/internal/channel"
	"github.com/5G-Measure-India/collector/internal/config.go"
	"github.com/5G-Measure-India/collector/internal/util"
)

var done channel.Type

type NdtFinalInfo struct {
	Download struct {
		Throughput Value `json:"Throughput"`
		Latency    Value `json:"Latency"`
	} `json:"Download"`
	Upload struct {
		Throughput Value `json:"Throughput"`
		Latency    Value `json:"Latency"`
	} `json:"Upload"`
}

type Value struct {
	Value float64 `json:"Value"`
}

func SpeedtestRoutine() {
	defer close(channel.SpeedtestDone)

	lastMinute := time.Now().Minute()

	csvFile := util.GetFilePath(config.OutDir, "speedtest.csv", config.Timestamp)
	csvWriter, err := os.Create(csvFile)
	if err != nil {
		log.Println("[speedtest] error opening log file:", err)
		return
	}
	defer csvWriter.Close()

	if _, err := csvWriter.WriteString("timestamp,dl_tp,dl_lat,ul_tp,ul_lat\n"); err != nil {
		log.Println("[speedtest] error writing csv header:", err)
		return
	}

	logIn, logOut := io.Pipe()
	defer logOut.Close()
	defer logIn.Close()

	for {
		select {
		case <-channel.Stop:
			return
		case <-time.After(5 * time.Second):
			if time.Now().Minute() == lastMinute {
				continue
			}

			lastMinute = time.Now().Minute()
			done = make(channel.Type)

			cmd := exec.Command("adb", "shell", "/data/local/tmp/ndt7-client", "-format", "json")
			cmd.Stdout = logOut

			if err := cmd.Start(); err != nil {
				log.Println("[speedtest] error starting:", err)
				return
			}
			log.Println("[speedtest] started | logging to:", csvFile)

			go logger(logIn, csvWriter)

			select {
			case <-channel.Stop:
				if err := cmd.Process.Kill(); err != nil {
					log.Println("[speedtest] error stopping:", err)
				}
				return
			case <-done:
				if _, err := cmd.Process.Wait(); err != nil {
					log.Println("[speedtest] error stopping:", err)
					return
				}
			}
		}
	}
}

func logger(logPipe *io.PipeReader, writer *os.File) {
	defer close(done)

	var (
		line    string
		ndtInfo NdtFinalInfo
	)

	scanner := bufio.NewScanner(logPipe)
	for scanner.Scan() {
		line = scanner.Text()

		if strings.Contains(line, "ServerFQDN") {
			break
		}
	}

	if json.Unmarshal([]byte(line), &ndtInfo) == nil {
		fmt.Fprintf(writer, "%f,%f,%f,%f,%f\n", util.GetTime(), ndtInfo.Download.Throughput.Value, ndtInfo.Download.Latency.Value, ndtInfo.Upload.Throughput.Value, ndtInfo.Upload.Latency.Value)
	}
}
