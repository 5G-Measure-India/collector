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
	close(util.GeoDone)

	csvFile := path.Join(config.DATA_DIR, config.GetFname("geo.log"))

	if err := config.Adb("shell", "am", "start", "-a", "com.llamalab.automate.intent.action.START_FLOW", "-d", fmt.Sprintf("content://com.llamalab.automate.provider/flows/%d/statements/1", LOG_FLOW_ID), "-e", "path", csvFile, "-e", "delay", "10").Run(); err != nil {
		log.Println("[geo] error starting flow:", err)
		return
	}
	log.Println("[phy] started | logging to: adb://", csvFile)

	if err := config.Adb("shell", "am", "start", "-a", "com.llamalab.automate.intent.action.START_FLOW", "-d", fmt.Sprintf("content://com.llamalab.automate.provider/flows/%d/statements/1", STOP_FLOW_ID), "-e", "flow", strconv.Itoa(LOG_FLOW_ID)).Run(); err != nil {
		log.Println("[geo] error stopping flow:", err)
	}

	<-util.Stop

	if err := config.Adb("pull", csvFile, config.OutDir).Run(); err != nil {
		log.Println("[geo] error pulling logs:", err)
	}
}