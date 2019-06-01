package main

import (
	"fmt"
	"regexp"
	"strings"
)

func createReply(userID string, message string) string {
	expression := strings.TrimPrefix(message, "!r ")
	ok := checkExpression(expression)
	if !ok {
		return fmt.Sprintf("<@!%s> I can't understand whatcha sayin', sorry", userID)
	}
	result, ok := roll(expression)
	if !ok {
		return fmt.Sprintf("<@!%s> I can't understand whatcha sayin', sorry", userID)
	}
	return fmt.Sprintf("<@!%s> `%s` = **%d**", userID, expression, result)
}

func checkExpression(expression string) bool {
	lang := regexp.MustCompile(`^(?:(?:\+|\-|^)(?:\d+|\d*d\d+(?:k\d+)?))+$`)
	return lang.MatchString(expression)
}

func roll(expression string) (int, bool) {
	atom := regexp.MustCompile(`(\+|\-|^)((\d*)d(\d+)(k\d+)?|(\d+))`)

	result := 0
	matches := atom.FindAllStringSubmatch(expression, -1)
	for _, throw := range matches {
		//handle single dice
		fmt.Println("Found dice: ", throw)
		result++
	}

	return result, true
}
