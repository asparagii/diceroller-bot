package main

import (
	"fmt"
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
		result, description, err := lang.Parse(expr)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("%s => %d\n", description, result)
		}
	}
}
