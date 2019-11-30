package main

import (
	"fmt"
	"github.com/MicheleLambertucci/diceroller-bot/pkg/lang"
	"strings"
)

func createReply(userID string, message string) string {
	expression := strings.TrimPrefix(message, "!r ")

	result, representation, err := lang.Start(expression)

	var response string

	if err != nil {
		response = fmt.Sprintf("<@!%s> Error: %v", userID, err)
	} else {
		response = fmt.Sprintf("<@!%s> `%s` => %s = `%v`", userID, expression, representation, result)

		if len(response) > 1800 {
			response = fmt.Sprintf("<@!%s> `%s` => `%d`", userID, expression, result)
		}
	}

	return response
}
