package main

import (
	"os"

	"lox/internal/core"

	"github.com/sirupsen/logrus"
)

func main() {
	core.RunFile("/Users/qianhangjian/interpreter/test/fib.lox")
	return

	if len(os.Args) > 2 {
		logrus.Errorln("invalid args, usage: ")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		// run file
		core.RunFile("/Users/qianhangjian/interpreter/test/token.lox")
	} else {
		// run prompt
		core.RunPrompt()
	}
}
