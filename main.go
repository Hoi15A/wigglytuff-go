package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"plugin"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token  string
	Prefix string
)

var commands map[string]func(string, *discordgo.Session, *discordgo.MessageCreate)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&Prefix, "p", ">>", "Bot Prefix")
	flag.Parse()

	// commands := make(map[string]func(string))
	// Load Commands
	commands = loadPlugins()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(m.Content, Prefix) {
		return
	}
	m.Content = strings.TrimPrefix(m.Content, Prefix)

	commandName := strings.Split(m.Content, " ")[0]
	args := strings.Replace(m.Content, commandName+" ", "", 1)
	if val, ok := commands[commandName]; ok {
		val(args, s, m)
	}
}

func loadPlugins() map[string]func(string, *discordgo.Session, *discordgo.MessageCreate) {
	commands := make(map[string]func(string, *discordgo.Session, *discordgo.MessageCreate))
	files, err := ioutil.ReadDir("./build")

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fileName := file.Name()
		fmt.Println("file:", fileName)

		if strings.HasSuffix(fileName, ".so") { // Note, dont keep this stupid extension
			p, err := plugin.Open("./build/" + fileName)

			if err != nil {
				panic(err)
			}
			regCommand, err := p.Lookup("RegisterCommand")
			if err != nil {
				panic(err)
			}
			_name, err := p.Lookup("Name")
			if err != nil {
				panic(err)
			}

			var name string = *(_name.(*string))

			commands[name] = regCommand.(func(string, *discordgo.Session, *discordgo.MessageCreate))

			fmt.Println("name:", name)
		}

	}
	return commands
}
