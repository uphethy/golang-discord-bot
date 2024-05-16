package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	sess, err := discordgo.New("Bot MTE3OTQ4MTk1MzcwODI4NjA4Mg.GcxHpR.rNOwD3zxh7wQpZ0M4yd7R0eGjMb4RSazEkdNFg")
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("mysql", "root:622gresko#@tcp(localhost:3306)/commands")
	if err != nil {
		fmt.Println("database error")
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	const passaway = "_-A2kCm#(*^)"
	go sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// guild, err := s.State.Guild(m.GuildID)
		// if err != nil {
		// 	return
		// }
		args := strings.Split(m.Content, " ")
		if m.Author.ID == s.State.User.ID {
			return
		}

		// isAdmin, err := isAdmin(s, m.Author.ID, m.GuildID)
		// if err != nil {
		//	return
		// }

		// vcommand zbase  && (isAdmin || m.Author.ID == guild.OwnerID)
		if args[0] == "vcommand" {
			if len(args) > 2 {
				s.ChannelMessageSend(m.ChannelID, "The command must consist of one word")
			} else if len(args) > 1 && args[1] != "" {
				insert, err := db.Query("INSERT INTO `commandContent` (`guild_id`, `command`, `content`) VALUES ('" + m.GuildID + "', '" + args[1] + "', '" + passaway + "')")
				if err != nil {
					log.Fatal(err)
				}
				err1 := s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
				if err1 != nil {
					log.Fatal(err1)
				}
				defer insert.Close()
			}
		} else if args[0] == "vhelp" {
			s.ChannelMessageSend(m.ChannelID, "vcommand <new command> - add a new command\n<command> add <some content> - add a content to command\n<command> remove <content> - remove a content from command")
		} else if len(args) == 1 && args[0] != "vcommand" { // zbase
			query := fmt.Sprintf("SELECT content FROM commandContent WHERE guild_id = '%s' AND command = '%s' ORDER BY RAND() LIMIT 1", m.GuildID, args[0])
			res, err := db.Query(query)
			if err != nil {
				log.Fatal(err)
			}
			defer res.Close()

			var content string
			if res.Next() {
				err := res.Scan(&content)
				if err != nil {
					log.Fatal(err)
				}
			}
			if content == passaway {
				return
			}
			s.ChannelMessageSend(m.ChannelID, content)
		} else if len(args) > 2 {
			if args[1] == "add" { // zbase add blablabla
				query := fmt.Sprintf("SELECT command FROM commandContent WHERE guild_id = '%s'", m.GuildID)
				res, err := db.Query(query)
				if err != nil {
					log.Fatal(err)
				}
				for res.Next() {
					var commanda string
					err = res.Scan(&commanda)
					if err != nil {
						log.Fatal(err)
					}
					defer res.Close()
					if commanda == args[0] {
						contenta := args[2:]
						insert, err := db.Query("INSERT INTO `commandContent` (`guild_id`, `command`, `content`) VALUES ('" + m.GuildID + "', '" + args[0] + "', '" + strings.Join(contenta, " ") + "')")
						if err != nil {
							log.Fatal(err)
						}
						defer insert.Close()
						queryx := fmt.Sprintf("DELETE FROM commandContent WHERE content = '%s'", passaway)
						_, err = db.Query(queryx)
						if err != nil {
							log.Fatal(err)
						}
						err1 := s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
						if err1 != nil {
							log.Fatal(err1)
						}
						break
					}
				}
			} else if args[1] == "remove" {
				query := fmt.Sprintf("SELECT content FROM commandContent WHERE guild_id = '%s' AND command = '%s'", m.GuildID, args[0])
				res, err := db.Query(query)
				if err != nil {
					log.Fatal(err)
				}
				for res.Next() {
					var content string
					err = res.Scan(&content)
					if err != nil {
						log.Fatal(err)
					}
					if content == strings.Join(args[2:], " ") {
						queryx := fmt.Sprintf("DELETE FROM commandContent WHERE content = '%s'", content)
						_, err = db.Query(queryx)
						if err != nil {
							log.Fatal(err)
						}
						err1 := s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
						if err1 != nil {
							log.Fatal(err1)
						}
						break
					}
				}
				err1 := s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
				if err1 != nil {
					log.Fatal(err1)
				}
			}

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

func isAdmin(s *discordgo.Session, userID string, guildID string) (bool, error) {
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		return false, err
	}

	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			continue
		}

		if role.Permissions&discordgo.PermissionAdministrator != 0 {
			return true, nil
		}
	}

	return false, nil
}
