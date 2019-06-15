package main

import (
	"fmt"
	"github.com/MicheleLambertucci/diceroller-bot/pkg/lang"
	"strings"
)

func createReply(userID string, message string) string {
	expression := strings.TrimPrefix(message, "!r ")

	result, err := lang.Parse(expression)

	var response string

	if err != nil {
		response = fmt.Sprintf("<@!%s> Error: %v", userID, err)
	} else {
		if result.Type() == lang.NUMBERVALUE {
			response = fmt.Sprintf("<@!%s> `%s` => %s = `%d`", userID, expression, result, result.(lang.Number).V())
		} else {
			response = fmt.Sprintf("<@!%s> `%s` => %s", userID, expression, result)
		}
		if len(response) > 1800 {
			response = fmt.Sprintf("<@!%s> `%s` => `%d`", userID, expression, result)
		}
	}

	return response
}
