package core

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func RunFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// TODO: implement
	_ = data
	return nil
}

func RunPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("> ")
		text, _ := reader.ReadString('\n')
		if len(text) == 0 {
			break
		}
		run(text)
	}
}

func run(source string) {
	scan := NewScanner(source)

	scan.scanTokens()

	for _, v := range scan.tokens {
		fmt.Printf("token: %v\n", v)
	}
}
