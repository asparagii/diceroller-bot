package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

type ExpressionResult []Dice

func (e ExpressionResult) String() string {
	var result strings.Builder
	for _, v := range e {
		result.WriteString(fmt.Sprintf("%v", Dice(v)))
	}
	ret := result.String()
	if ret[0] == '+' {
		return ret[1:]
	}
	return ret
}

func (e ExpressionResult) total() int {
	sum := 0
	for _, dice := range e {
		sum += dice.sign() * dice.total()
	}
	return sum
}

type Dice struct {
	number   int
	size     int
	negative bool
	keep     int
	scores   []int
}

func (d Dice) sign() int {
	if d.negative {
		return -1
	} else {
		return 1
	}
}

func (d Dice) String() string {
	var signRepr string
	if d.negative {
		signRepr = "-"
	} else {
		signRepr = "+"
	}

	if d.size == 1 {
		return fmt.Sprintf("%s%v", signRepr, d.scores[0])
	}

	var out []string
	for i, v := range d.scores {
		var highlight string
		if v == d.size {
			highlight = "**"
		} else if d.keep != -1 && i < len(d.scores)-d.keep {
			highlight = "~~"
		} else {
			highlight = ""
		}

		throwRepr := fmt.Sprintf("%s%v%s", highlight, v, highlight)
		out = append(out, throwRepr)
	}
	body := strings.Join(out, "+")
	return fmt.Sprintf("%s(%s)", signRepr, body)
}

func (d Dice) total() int {
	sum := 0
	for i, v := range d.scores {
		if d.keep != -1 {
			if i >= len(d.scores)-d.keep {
				sum += v
			}
		} else {
			sum += v
		}
	}
	return sum
}

func (d *Dice) extract() {
	if d.size == 1 {
		d.scores = make([]int, 1, 1)
		d.scores[0] = d.number
	} else {
		d.scores = make([]int, d.number, d.number)
		for i := 0; i < d.number; i++ {
			d.scores[i] = rand.Intn(d.size) + 1
		}
		sort.Sort(sort.IntSlice(d.scores))
	}
}
