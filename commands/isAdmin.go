package commands

import (
	"github.com/bwmarrin/discordgo"
)

func IsAdmin(s *discordgo.Session, userID string, guildID string) bool {
	member, _ := s.GuildMember(guildID, userID)

	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			continue
		}

		if role.Permissions&discordgo.PermissionAdministrator != 0 {
			return true
		}
	}
	return false
}
