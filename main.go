package main

import (
	"auto-initail-server/helper"
	"auto-initail-server/utils"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {

	// log config
	logName := "app.log"
	logFile, err := os.OpenFile(logName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// pars flag
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags]\n\tExample: auto-initail-server -c conf.yaml -f ~/.ssh/rsa_pub\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	ymlFile := flag.String("c", "", "yaml file path")
	PublicKey := flag.String("f", "", "PuplicKey")

	flag.Parse()

	if flag.NFlag() != 2 {
		flag.Usage()
		return
	}

	conf := utils.Newconfig(*ymlFile)
	if err := conf.GetConf(); err != nil {
		log.Fatal(err)
	}

	absPuplicKey, err := filepath.Abs(*PublicKey)
	if err != nil {
		log.Fatal(err)
	}

	sshObj := helper.NewInfoSSH(conf)
	if err := sshObj.RunQueue(absPuplicKey); err != nil {
		log.Fatal(err)
	}

}
