package main

import (
	"flag"
	"fmt"
	"github.com/asparagii/diceroller-bot/lang"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

var token string

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	if token == "" {
		fmt.Println("No token provided. Please run: diceroller -t <bot-token>")
		return
	}

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session: ", err)
		return
	}

	discord.AddHandler(ready)
	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening Discord session: ", err)
	}

	fmt.Println("Diceroller is now running. Send SIGINT or SIGTERM to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, "Fortnite")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!r ") {
		response := createReply(m.Author.ID, m.Content)
		mess, err := s.ChannelMessageSend(m.ChannelID, response)
		if err != nil {
			fmt.Println("Unable to send message to channel: ", err)
		} else {
			fmt.Println("Sent message: ", mess.ID, mess.Content)
		}
	}
}

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
