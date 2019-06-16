package main

import (
	"fmt"
	"github.com/MicheleLambertucci/diceroller-bot/pkg/lang"
	"github.com/mattn/go-isatty"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	var expr string
	if isatty.IsTerminal(os.Stdin.Fd()) {
		fmt.Println("Welcome to dicelang interactive prompt")
		fmt.Println("Type 'exit' or 'quit' to exit")
		for true {
			fmt.Print("> ")
			fmt.Scanf("%s", &expr)
			if strings.Compare(expr, "quit") == 0 || strings.Compare(expr, "exit") == 0 {
				break
			}
			compute(expr)
		}
	} else {
		for true {
			_, err := fmt.Scanf("%s", &expr)
			if err != nil {
				break
			}
			compute(expr)
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
	result, err := lang.Parse(expr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		//	if result.Properties[lang.TYPE] == lang.NUMBERVALUE {
		//		fmt.Printf("%v => %d\n", cliFormatter(result.String()), result.(lang.Number).V())
		//	} else {
		fmt.Printf("%v\n", result)
		//}
	}
}
