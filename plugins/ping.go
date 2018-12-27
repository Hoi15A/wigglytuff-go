package main

import (
	"github.com/bwmarrin/discordgo"
)

var Name string = "ping"

func RegisterCommand(test string, s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Pong!")
}
