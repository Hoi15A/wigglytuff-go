package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var Name string = "getroles"

func RegisterCommand(test string, s *discordgo.Session, m *discordgo.MessageCreate) {
	var guild, err = s.Guild(m.GuildID)
	if err != nil {
		fmt.Println("Guild "+m.GuildID+" not found: ", err)
		return
	}

	var response string = "**__Roles:__** \n"
	for _, role := range guild.Roles {
		if role.Name == "@everyone" {
			continue // cuz for some reason everyone is a role, thanks discord
		}
		response += role.Name + ": " + role.ID + "\n"
	}

	s.ChannelMessageSend(m.ChannelID, response)
}
