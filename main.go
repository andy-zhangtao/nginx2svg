package main

import (
	"github.com/sirupsen/logrus"
)

const modulename = "ng2svg"

var version = "0.1.0"

func main() {
	logrus.WithFields(logrus.Fields{"version": version}).Info(modulename)
	start()
}
