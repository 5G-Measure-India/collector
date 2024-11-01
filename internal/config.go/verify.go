package config

import (
	"errors"
	"log"
	"os"
)

func checkOutDir() error {
	if fd, err := os.Stat(OutDir); err != nil {
		if err := os.MkdirAll(OutDir, 0755); err != nil {
			return errors.New("could not create directory: " + err.Error())
		}
	} else if !fd.IsDir() {
		return errors.New("not a directory")
	}

	if tmpFile, err := os.CreateTemp(OutDir, ".write"); err != nil {
		return errors.New("requires write access")
	} else {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}

	return nil
}

func verifyFlags() {
	if err := checkOutDir(); err != nil {
		log.Println("[config] output dir:", OutDir, err)
		os.Exit(1)
	} else {
		log.Println("[config] output dir:", OutDir)
	}
}
