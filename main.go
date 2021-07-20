package main

import (
	"github.com/programmablemike/fibo/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.InfoLevel)
	cmd.Execute()
}
