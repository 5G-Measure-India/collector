package net

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/5G-Measure-India/collector/internal/config"
	"github.com/5G-Measure-India/collector/internal/util"
)

var done util.Type

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
	defer close(util.SpeedtestDone)

	lastMinute := time.Now().Minute()

	csvFile := path.Join(config.OutDir, config.GetFname("speedtest.csv"))
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
	log.Println("[speedtest] started | logging to:", csvFile)

	logIn, logOut := io.Pipe()
	defer logOut.Close()
	defer logIn.Close()

	for {
		select {
		case <-util.Stop:
			return
		case <-time.After(5 * time.Second):
			if time.Now().Minute() == lastMinute {
				continue
			}

			lastMinute = time.Now().Minute()
			done = make(util.Type)

			cmd := config.Adb("shell", path.Join(config.TOOLS_DIR, "ndt7-client"), "-format", "json")
			cmd.Stdout = logOut

			if err := cmd.Start(); err != nil {
				log.Println("[speedtest] error starting:", err)
				return
			}

			go logger(logIn, csvWriter)

			select {
			case <-util.Stop:
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
