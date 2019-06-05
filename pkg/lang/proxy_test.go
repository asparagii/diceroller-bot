package lang

import (
	"testing"
)

func TestProxy(t *testing.T) {
	channel := make(chan Token)
	proxy := Nexter{token: channel}
	go Inputs(channel)
	proxy.Init()
	Expect(t, proxy.Peek(), Token{t: NUMBER})
	Expect(t, proxy.Pop(), Token{t: NUMBER})
	Expect(t, proxy.Peek(), Token{t: EOF})
	Expect(t, proxy.Pop(), Token{t: EOF})
	Expect(t, proxy.Peek(), Token{t: EOF})
	Expect(t, proxy.Pop(), Token{t: EOF})
}

func Inputs(channel chan Token) {
	channel <- Token{t: NUMBER}
	channel <- Token{t: EOF}
}
