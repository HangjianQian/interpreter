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

	run(string(data))
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

	// for _, v := range scan.tokens {
	// 	fmt.Printf("token: %+v\n", v)
	// }

	parser := NewParser(scan.tokens)
	stmts := parser.doParse()
	// fmt.Printf("stmt s%+v\n", stmts)

	inter := NewInterpreter()

	resolver := NewResolver(inter)
	resolver.resolve(BlockStmt{stmts: stmts})

	for _, s := range stmts {
		inter.interpret(s)
		// fmt.Printf("res: %v\n", inter.interpret(s))
	}
}
