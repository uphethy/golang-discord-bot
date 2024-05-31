package main

import (
	"database/sql"
	"fmt"
	"golang-discord-bot/commands"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	sql_connection := os.Getenv("SQL_CONNECTION")

	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("bot token error")
	}
	db, err := sql.Open("mysql", sql_connection)
	if err != nil {
		fmt.Println("database error")
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		args := strings.Split(m.Content, " ")

		if m.Author.ID == s.State.User.ID {
			return
		}

		if args[0] == "vcommand" {
			if len(args) > 2 {
				s.ChannelMessageSend(m.ChannelID, "the command must consist of one word")
			} else if len(args) > 1 && args[1] != "" {
				go commands.AddCommand(s, m, db, args)
			}

		} else if args[0] == "vhelp" {
			s.ChannelMessageSend(m.ChannelID, "```vcommand <new command> - add a new command\n<command> add <some content> - add a content to command\n<command> remove <content> - remove a content from command```")
		} else if len(args) > 1 && args[1] == "add" && args[0] != "vcommand" && args[0] != "vhelp" {
			if len(args) == 2 {
				err := s.MessageReactionAdd(m.ChannelID, m.ID, "❎")
				if err != nil {
					panic(err)
				}
				s.ChannelMessageSend(m.ChannelID, "empty command")
			} else {
				go commands.AddCommandContent(s, m, db, args)
			}
		} else if len(args) > 1 && args[1] == "remove" && args[0] != "vcommand" && args[0] != "vhelp" {
			if len(args) == 2 {
				err := s.MessageReactionAdd(m.ChannelID, m.ID, "❎")
				if err != nil {
					panic(err)
				}
				s.ChannelMessageSend(m.ChannelID, "nothing to remove")
			} else {
				go commands.RemoveContent(s, m, db, args)
			}
		} else if args[0] == "vfor" && len(args) == 3 {
			count, err := strconv.Atoi(args[2])
			if err != nil || count > 10 {
				s.MessageReactionAdd(m.ChannelID, m.ID, "❎")
				return
			}
			for i := 0; i < count; i++ {
				commands.SendRandomContentFor(s, m, db, args)
			}
		} else if len(args) >= 1 && args[0] != "vcommand" && args[0] != "vhelp" {
			go commands.SendRandomContent(s, m, db, args)
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
