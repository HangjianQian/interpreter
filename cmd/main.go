package main

import (
	"os"

	"lox/internal/core"

	"github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) > 2 {
		logrus.Errorln("invalid args, usage: ")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		// run file

	} else {
		// run prompt
		core.RunPrompt()
	}
}
