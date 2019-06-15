package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/MicheleLambertucci/diceroller-bot/pkg/lang"
)

func main() {
	var expr string
	for true {
		fmt.Print("> ")
		fmt.Scanf("%s", &expr)
		if strings.Compare(expr, "quit") == 0 {
			break
		}
		result, err := lang.Parse(expr)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			if result.Type() == lang.NUMBERVALUE {
				fmt.Printf("%v => %d\n", cliFormatter(result.String()), result.(lang.Number).V())
			} else {
				fmt.Printf("%v\n", result)
			}
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
