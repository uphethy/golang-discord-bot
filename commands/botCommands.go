package commands

import (
	"database/sql"
	"fmt"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
)

func AddCommand(s *discordgo.Session, m *discordgo.MessageCreate, db *sql.DB, args []string) {
	insert, err := db.Query("INSERT INTO `commands` (`guild_ID`, `command`) VALUES ('" + m.GuildID + "', '" + args[1] + "')")
	if err != nil {
		err1 := s.MessageReactionAdd(m.ChannelID, m.ID, "❎")
		if err1 != nil {
			panic(err1)
		}
		s.ChannelMessageSend(m.ChannelID, "that command already exists")
		return
	}
	err1 := s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	if err1 != nil {
		panic(err1)
	}
	defer insert.Close()
}

func AddCommandContent(s *discordgo.Session, m *discordgo.MessageCreate, db *sql.DB, args []string) {
	var id_command int64
	err := db.QueryRow("SELECT id FROM commands WHERE guild_ID = '" + m.GuildID + "' AND command = '" + args[0] + "'").Scan(&id_command)
	if err != nil {
		panic(err)
	}
	query := fmt.Sprintf("INSERT INTO `command_contents` (`command_ID`, `content`) VALUES ('%d', '%s')", id_command, args[2])
	insert, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	err1 := s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	if err1 != nil {
		panic(err1)
	}
	defer insert.Close()
}

func SendRandomContent(s *discordgo.Session, m *discordgo.MessageCreate, db *sql.DB, args []string) {
	var id_command int64
	err := db.QueryRow("SELECT id FROM commands WHERE guild_ID = '" + m.GuildID + "' AND command = '" + args[0] + "'").Scan(&id_command)
	if err != nil {
		return
	}
	var content string
	query := fmt.Sprintf("SELECT content FROM command_contents WHERE command_id = '%d' ORDER BY RAND() LIMIT 1", id_command)
	err = db.QueryRow(query).Scan(&content)
	if err != nil {
		return
	}
	s.ChannelMessageSend(m.ChannelID, content)
}
