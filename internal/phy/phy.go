package phy

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/5G-Measure-India/collector/internal/config"
	"github.com/5G-Measure-India/collector/internal/util"
)

var done = make(util.Type)

type Ml1Info struct {
	Timestamp            string `json:"timestamp"`
	ComponentCarrierList []struct {
		Cells []struct {
			CellQualityRsrp string `json:"Cell Quality Rsrp"`
			CellQualityRsrq string `json:"Cell Quality Rsrq"`
		} `json:"Cells"`
	} `json:"Component_Carrier List"`
}

func unmarshalTime(data string) (tstamp time.Time, err error) {
	tstamp, err = time.Parse("2006-01-02 15:04:05.000000", data)
	if err == nil {
		tstamp = tstamp.In(time.Local)
	}

	return
}

func PhyRoutine() {
	defer close(util.PhyDone)

	csvFile := path.Join(config.OutDir, config.GetFname("phy.log"))
	csvWriter, err := os.Create(csvFile)
	if err != nil {
		log.Println("[phy] error opening log file:", err)
		return
	}
	defer csvWriter.Close()

	if _, err := csvWriter.WriteString("timestamp,rsrp,rsrq\n"); err != nil {
		log.Println("[phy] error writing csv header:", err)
		return
	}

	logIn, logOut := io.Pipe()
	defer logOut.Close()
	defer logIn.Close()

	cmd := exec.Command(config.Python, config.MobMon)
	cmd.Stdout = logOut
	cmd.Stderr = logOut

	if err := cmd.Start(); err != nil {
		log.Println("[phy] error starting:", err)
		return
	}
	log.Println("[phy] started | logging to:", csvFile)

	go logger(logIn, csvWriter)

	<-util.Stop

	if err := cmd.Process.Kill(); err != nil {
		log.Println("[phy] error stopping:", err)
	} else if _, err := cmd.Process.Wait(); err != nil {
		log.Println("[phy] error stopping:", err)
	}

	<-done
}

func logger(logPipe *io.PipeReader, writer *os.File) {
	defer close(done)

	var (
		line string
		// ml1Info Ml1Info
	)

	scanner := bufio.NewScanner(logPipe)
	for scanner.Scan() {
		line = scanner.Text()

		if strings.Contains(line, "5G_NR_ML1_Searcher_Measurement_Database_Update_Ext") {
			if i := strings.Index(line, "{"); i != -1 {
				// if json.Unmarshal([]byte(line[i:]), &ml1Info) == nil {
				// 	if tstamp, err := unmarshalTime(ml1Info.Timestamp); err == nil {
				// 		fmt.Fprintf(writer, "%f,%s,%s,0\n", util.GetTime(tstamp), ml1Info.ComponentCarrierList[0].Cells[0].CellQualityRsrp, ml1Info.ComponentCarrierList[0].Cells[0].CellQualityRsrq)
				// 	}
				// }
				fmt.Fprintln(writer, line[i:])
			}
		}
	}
}
