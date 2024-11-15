package phy

import (
	"fmt"
	"log"
	"path"
	"strconv"

	"github.com/5G-Measure-India/collector/internal/config"
	"github.com/5G-Measure-India/collector/internal/util"
)

const (
	LOG_FLOW_ID  = 7
	STOP_FLOW_ID = 8
)

func GeoRoutine() {
	defer close(util.GeoDone)

	fname := config.GetFname("geo.log")
	csvFile := path.Join(config.DATA_DIR, fname)
	dataFile := path.Join(config.OutDir, fname)

	if err := config.Adb("shell", "am", "start", "-a", "com.llamalab.automate.intent.action.START_FLOW", "-d", fmt.Sprintf("content://com.llamalab.automate.provider/flows/%d/statements/1", LOG_FLOW_ID), "-e", "path", csvFile, "-e", "delay", "10").Run(); err != nil {
		log.Println("[geo] error starting flow:", err)
		return
	}
	log.Println("[phy] started | logging to:", "adb://"+csvFile)

	<-util.Stop

	if err := config.Adb("shell", "am", "start", "-a", "com.llamalab.automate.intent.action.START_FLOW", "-d", fmt.Sprintf("content://com.llamalab.automate.provider/flows/%d/statements/1", STOP_FLOW_ID), "-e", "flow", strconv.Itoa(LOG_FLOW_ID)).Run(); err != nil {
		log.Println("[geo] error stopping flow:", err)
	}

	if err := config.Adb("pull", csvFile, dataFile).Run(); err != nil {
		log.Println("[geo] error pulling logs:", err)
	} else {
		log.Println("[geo] logs pulled:", dataFile)
	}

	config.Adb("shell", "rm", csvFile).Run()
}
