package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/asparagii/diceroller-bot/lang"
	"github.com/mattn/go-isatty"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) > 1 {
		// process single input from args and exit
		inputProgram := os.Args[1]
		compute(inputProgram)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	var inputProgram string
	if isatty.IsTerminal(os.Stdin.Fd()) {
		fmt.Println("Welcome to dicelang interactive prompt")
		fmt.Println("Type 'exit' or 'quit' to exit")
		for {
			fmt.Print("> ")
			scanner.Scan()
			inputProgram = scanner.Text()

			if strings.Compare(inputProgram, "quit") == 0 || strings.Compare(inputProgram, "exit") == 0 {
				break
			}
			compute(inputProgram)
		}
	} else {
		for {
			scanner.Scan()
			inputProgram = scanner.Text()
			compute(inputProgram)
		}
	}
}

func cliFormatter(str string) string {
	strong := regexp.MustCompile(`\*\*(.+?)\*\*`)
	discard := regexp.MustCompile(`\~\~(.+?)\~\~`)

	ret := strong.ReplaceAllString(str, "\033[1;33m$1\033[0m")
	ret = discard.ReplaceAllString(ret, "\033[9m$1\033[0m")

	return ret
}

func compute(expr string) {
	result, representation, err := lang.Start(expr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("%s => %v\n", cliFormatter(representation), result)
	}
}
