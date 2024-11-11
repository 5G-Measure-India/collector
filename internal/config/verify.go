package config

import (
	"errors"
	"log"
	"os"
	"os/exec"
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

func checkAdb() error {
	if exec.Command(adb, "--version").Run() != nil {
		return errors.New("not installed")
	}

	return nil
}

func checkPython() error {
	if exec.Command(Python, "--version").Run() != nil {
		return errors.New("not installed")
	}

	return nil
}

func checkMobMon() error {
	return nil
}

func checkPingServer() error {
	return nil
}

func verifyFlags() {
	if err := checkOutDir(); err != nil {
		log.Println("[config] output dir:", OutDir, err)
		os.Exit(1)
	} else {
		log.Println("[config] output dir:", OutDir)
	}

	if err := checkAdb(); err != nil {
		log.Println("[config] adb:", adb, err)
		os.Exit(1)
	} else {
		log.Println("[config] adb:", adb)
	}

	if err := checkPython(); err != nil {
		log.Println("[config] python:", Python, err)
		os.Exit(1)
	} else {
		log.Println("[config] python:", Python)
	}

	if err := checkMobMon(); err != nil {
		log.Println("[config] mobile-insight monitor:", MobMon, err)
		os.Exit(1)
	} else {
		log.Println("[config] mobile-insight monitor:", MobMon)
	}

	if err := checkPingServer(); err != nil {
		log.Println("[config] ping server:", PingServer, err)
		os.Exit(1)
	} else {
		log.Println("[config] ping server:", PingServer)
	}
}
