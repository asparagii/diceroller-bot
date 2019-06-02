package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var InvalidExpression = errors.New("Expression is invalid")

func createReply(userID string, message string) (string, error) {
	expression := strings.TrimPrefix(message, "!r ")
	expression = StripWhiteSpace(expression)
	ok := checkExpression(expression)
	if !ok {
		return "", InvalidExpression
	}
	result, ok := roll(expression)
	if !ok {
		return "", InvalidExpression
	}
	resultMessage := fmt.Sprintf("<@!%s> %v = `%d`", userID, result, result.total())
	if len(resultMessage) > 1900 {
		return fmt.Sprintf("<@!%s> (...) = `%d`", userID, result.total()), nil
	}
	return resultMessage, nil
}

func checkExpression(expression string) bool {
	lang := regexp.MustCompile(`^(?:(?:\+|\-|^)(?:\d+|\d*d\d+(?:k\d+)?))+$`)
	return lang.MatchString(expression)
}

func roll(expression string) (ExpressionResult, bool) {
	atomExpr := regexp.MustCompile(`(\+|\-|^)(\d*d\d+(k\d+)?|\d+)`)
	matches := atomExpr.FindAllString(expression, -1)
	var result ExpressionResult
	for _, expr := range matches {
		dice, ok := createDice(expr)
		if !ok {
			fmt.Errorf("Error while creating dice for %s", expr)
			return result, false
		}
		dice.extract()
		result = append(result, dice)
	}
	return result, true
}

func StripWhiteSpace(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func createDice(atom string) (Dice, bool) {
	var ret Dice

	commonDice := regexp.MustCompile(`^(\+|\-)?(\d*)(d\d+)(k\d+)?$`)
	scalar := regexp.MustCompile(`^(\+|\-)?(\d+)$`)

	if scalar.MatchString(atom) {
		return createScalar(atom)
	} else if commonDice.MatchString(atom) {
		return createCommonDice(atom)
	} else {
		return ret, false
	}
}

func createScalar(expression string) (Dice, bool) {
	scalar := regexp.MustCompile(`^(\+|\-)?(\d+)$`)
	matches := scalar.FindStringSubmatch(expression)

	var ret Dice
	ret.negative = strings.Compare(matches[1], "-") == 0
	if len(matches[2]) == 0 {
		return ret, false
	}

	number, err := strconv.Atoi(matches[2])
	if err != nil {
		fmt.Errorf("Error while parsing Scalar expression '%s': %v", matches[2], err)
		return ret, false
	}
	ret.number = number
	ret.size = 1
	ret.keep = -1

	return ret, true
}

func createCommonDice(expression string) (Dice, bool) {
	commonDice := regexp.MustCompile(`^(\+|\-)?(\d*)(d\d+)(k\d+)?$`)
	matches := commonDice.FindStringSubmatch(expression)

	var ret Dice
	ret.negative = strings.Compare(matches[1], "-") == 0

	if len(matches[2]) > 0 {
		number, err := strconv.Atoi(matches[2])
		if err != nil {
			fmt.Errorf("Error while parsing CommonDice '%s': %v", matches[2], err)
			return ret, false
		}
		if number > 20000 {
			fmt.Errorf("Received `number` %d, which is bigger than maximum %d", number, 20000)
			return ret, false
		}
		ret.number = number
	} else {
		ret.number = 1
	}
	if len(matches[3]) > 0 {
		size, err := strconv.Atoi(matches[3][1:])
		if err != nil {
			return ret, false
		}
		ret.size = size
	}
	if len(matches[4]) > 0 {
		keep, err := strconv.Atoi(matches[4][1:])
		if err != nil {
			fmt.Errorf("Error while parsing CommonDice `keep` '%s': %v", matches[4][1:], err)
			return ret, false
		}
		if keep > ret.number {
			fmt.Errorf("Received `keep` bigger than `number`: %d > %d", keep, ret.number)
			return ret, false
		}
		ret.keep = keep
	} else {
		ret.keep = -1
	}
	return ret, true
}
