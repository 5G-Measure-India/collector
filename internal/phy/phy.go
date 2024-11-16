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

	"github.com/5G-Measure-India/collector/internal/config"
	"github.com/5G-Measure-India/collector/internal/util"
)

const miMonFile = "mi.py"

var done = make(util.Type)

// type Ml1SearcherData struct {
// 	Timestamp              string `json:"timestamp"`
// 	NumLayers              int32  `json:"Num Layers"`
// 	SSBPeriodicityServCell int32  `json:"SSB Periodicity Serv Cell"`
// 	FreqOffset             string `json:"Frequency Offset"`
// 	ComponentCarrierList   []struct {
// 		RasterARFCN      int32  `json:"Raster ARFCN"`
// 		ServingCellIndex int32  `json:"Serving Cell Index"`
// 		NumCells         int32  `json:"Num Cells"`
// 		ServingCellPCI   int32  `json:"Serving Cell PCI"`
// 		ServingSSB       int32  `json:"Serving SSB"`
// 		ServingRsrpRx230 string `json:"ServingRsrpRx23[0]"`
// 		ServingRsrpRx231 string `json:"ServingRsrpRx23[1]"`
// 		Cells            []Cell `json:"Cells"`
// 	} `json:"Component_Carrier List"`
// }

// type Cell struct {
// 	PCI             int32            `json:"PCI"`
// 	PBCHSFN         int32            `json:"PBCH SFN"`
// 	NumBeams        int32            `json:"Num Beams"`
// 	CellQualityRsrp string           `json:"Cell Quality Rsrp"`
// 	CellQualityRsrq string           `json:"Cell Quality Rsrq"`
// 	Beams           map[int]BeamInfo `json:"-"`
// }

// type BeamInfo struct {
// 	SSBIndex                int32  `json:"SSB Index"`
// 	RxBeamInfoRSRPs0        string `json:"RX Beam Info-RSRPs[0]"`
// 	RxBeamInfoRSRPs1        string `json:"RX Beam Info-RSRPs[1]"`
// 	Nr2NrFilteredBeamRsrpL3 string `json:"Nr2NrFilteredBeamRsrpL3"`
// 	Nr2NrFilteredBeamRsrqL3 string `json:"Nr2NrFilteredBeamRsrqL3"`
// }

// func unmarshalTime(data string) (tstamp time.Time, err error) {
// 	tstamp, err = time.Parse("2006-01-02 15:04:05.000000", data)
// 	if err == nil {
// 		tstamp = tstamp.In(time.Local)
// 	}

// 	return
// }

func PhyRoutine() {
	defer close(util.PhyDone)

	csvFile := path.Join(config.OutDir, config.GetFname("phy.log"))
	csvWriter, err := os.Create(csvFile)
	if err != nil {
		log.Println("[phy] error opening log file:", err)
		return
	}
	defer csvWriter.Close()

	// if _, err := csvWriter.WriteString("timestamp,rsrp,rsrq\n"); err != nil {
	// 	log.Println("[phy] error writing csv header:", err)
	// 	return
	// }

	logIn, logOut := io.Pipe()
	defer logOut.Close()
	defer logIn.Close()

	cmd := exec.Command(config.Python, miMonFile)
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
		// ml1SData Ml1SearcherData
	)

	scanner := bufio.NewScanner(logPipe)
	for scanner.Scan() {
		line = scanner.Text()

		if strings.Contains(line, "5G_NR_ML1_Searcher_Measurement_Database_Update_Ext") || strings.Contains(line, "5G_NR_RRC_OTA_Packet") {
			if i := strings.Index(line, "{"); i != -1 {
				fmt.Fprintln(writer, line[i:])
			}
		}
	}
}
