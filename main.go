/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	logrusStack "github.com/Gurpartap/logrus-stack"
	log "github.com/sirupsen/logrus"
	"io"
	"mystake/cmd"
	"mystake/lib"
	"time"
)

func init() {
	log.SetLevel(log.DebugLevel)
	callerLevels := []log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
	}
	stackLevels := []log.Level{log.PanicLevel, log.FatalLevel, log.ErrorLevel}
	log.AddHook(logrusStack.NewHook(callerLevels, stackLevels))
	log.AddHook(lib.RotateLogHook("log", "stdout.log", 7*24*time.Hour, 24*time.Hour))
	log.SetOutput(io.Discard)
}

func main() {
	cmd.Execute()
}
