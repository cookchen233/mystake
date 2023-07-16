/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"encoding/gob"
	"fmt"
	logrusStack "github.com/Gurpartap/logrus-stack"
	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"io"
	"mystake/cmd"
	"mystake/lib"
	"os"
	"time"
)

var ca *cache.Cache

func init() {
	er := godotenv.Load("./pro.env")
	if er != nil {
		fmt.Println(er)
		log.Error(er)
	}
	fmt.Println(os.Getenv("DIND_BOT_TOKEN"))
	cache_file := "./cache.gob"
	_, err := os.Lstat(cache_file)
	var M map[string]cache.Item
	if !os.IsNotExist(err) {
		File, _ := os.Open(cache_file)
		D := gob.NewDecoder(File)
		D.Decode(&M)
	}
	if len(M) > 0 {
		ca = cache.NewFrom(cache.NoExpiration, 10*time.Minute, M)
	} else {
		ca = cache.New(cache.NoExpiration, 10*time.Minute)
	}
	go func() {
		for {
			time.Sleep(time.Duration(60) * time.Second)
			File, _ := os.OpenFile(cache_file, os.O_RDWR|os.O_CREATE, 0777)
			defer File.Close()
			enc := gob.NewEncoder(File)
			if err := enc.Encode(ca.Items()); err != nil {
				log.Error(err)
			}
		}
	}()

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
