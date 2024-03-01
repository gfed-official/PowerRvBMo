package main

import (
	"fmt"
	"regexp"
    "log"

	"github.com/bwmarrin/discordgo"
)

func RevertHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channel, _ := s.Channel(i.ChannelID)
	category, _ := s.Channel(channel.ParentID)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Getting VMs...",
		},
	})

	r, err := regexp.Compile(`\d{4}$`)
	if err != nil {
		errMsg := fmt.Sprintf("Error getting team number: %v", err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

	teamNumber := r.Find([]byte(category.Name))
	if teamNumber == nil {
		errMsg := fmt.Sprintf("Error getting team number: %v", err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

	vms, err := GetVMs(string(teamNumber))
	if err != nil {
		errMsg := fmt.Sprintf("Error getting VMs: %v", err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

	msgArray := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    "revert-vm",
					Placeholder: "Select a VM",
					Options:     buildSelectOptions(vms),
				},
			},
		},
	}

	contentMsg := "Select a VM to revert to snapshot:"
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content:    &contentMsg,
		Components: &msgArray,
	})
}

func buildSelectOptions(vmList []string) []discordgo.SelectMenuOption {
	var options []discordgo.SelectMenuOption
	for _, vm := range vmList {
		options = append(options, discordgo.SelectMenuOption{
			Label:       vm,
			Value:       vm,
			Description: "Revert to snapshot",
			Emoji: discordgo.ComponentEmoji{
				Name: "🔧",
			},
		})
	}
	return options
}

func RevertVMSelectHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "Reverting VM...",
		},
	})

	err := Revert(data.Values[0])
	if err != nil {
		errMsg := fmt.Sprintf("Error reverting VM: %v", err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

	success := "VM reverted successfully!"

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &success,
	})
}

func PingHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	})
}

func TeamsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channel, _ := s.Channel(i.ChannelID)
	guild, _ := s.Guild(channel.GuildID)
	options := i.ApplicationCommandData().Options
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	switch options[0].Name {
	case "create":
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
		for _, opt := range options[0].Options {
			optionMap[opt.Name] = opt
		}
		id := int(optionMap["team-id"].IntValue())
		count := int(optionMap["team-count"].IntValue())
		for x := 0; x < count; x++ {
			createTeam(s, i, id+x, *guild)
		}
	case "delete":
		options = options[0].Options
		if options[0] == nil {
			errMsg := "Missing options"
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &errMsg,
			})
			return
		}
		switch options[0].Name {
		case "by-role":
			options = options[0].Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			role := optionMap["team-role"].RoleValue(s, guild.ID)
			deleteTeam(s, i, role, guild)
		case "all":
			groles, _ := s.GuildRoles(guild.ID)

			r, err := regexp.Compile(`\d{4}$`)
			if err != nil {
				errMsg := fmt.Sprintf("Error getting team number: %v", err)
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &errMsg,
				})
				return
			}

			for _, role := range groles {
				teamNumber := r.Find([]byte(role.Name))
				if teamNumber != nil {
					deleteTeam(s, i, role, guild)
				}
			}
		}
	}
}

func createTeam(s *discordgo.Session, i *discordgo.InteractionCreate, id int, guild discordgo.Guild) {
	teamName := fmt.Sprintf("Team %d", id)

	// Create team role
	msg := "Creating role: " + teamName
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})

    var blue int = 3447003
    var permissions int64 = 0
    var hoist bool = true
    var mentionable bool = true
    roleParam := discordgo.RoleParams{
        Name: teamName,
        Permissions: &permissions,
        Color: &blue, // blue
        Hoist: &hoist,
        Mentionable: &mentionable,
    }

    role, _ := s.GuildRoleCreate(guild.ID, &roleParam)

	msg = "Creating category: " + teamName
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
	// Create parent category
	category, _ := s.GuildChannelCreateComplex(guild.ID, discordgo.GuildChannelCreateData{
		Name: teamName,
		Type: 4,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:   guild.ID,
				Type: 0,
				Deny: discordgo.PermissionViewChannel,
			},
			{
				ID:    role.ID,
				Type:  0,
				Allow: discordgo.PermissionViewChannel,
			},
		},
	})

	// Create child channels
	msg = "Creating channels: " + teamName
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
	s.GuildChannelCreateComplex(guild.ID, discordgo.GuildChannelCreateData{
		Name:     fmt.Sprintf("%s-text", teamName),
		Type:     0,
		ParentID: category.ID,
	})
	s.GuildChannelCreateComplex(guild.ID, discordgo.GuildChannelCreateData{
		Name:     fmt.Sprintf("%s-support", teamName),
		Type:     0,
		ParentID: category.ID,
	})
	s.GuildChannelCreateComplex(guild.ID, discordgo.GuildChannelCreateData{
		Name:     fmt.Sprintf("%s-voice", teamName),
		Type:     2,
		ParentID: category.ID,
	})
}

func deleteTeam(s *discordgo.Session, i *discordgo.InteractionCreate, role *discordgo.Role, guild *discordgo.Guild) {
	roles, _ := s.GuildRoles(guild.ID)
	channels, _ := s.GuildChannels(guild.ID)
	parent := findChannelByName(s, i, role.Name)
    teamName := role.Name

	msg := "Deleting channels: " + teamName
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
	for _, gchannel := range channels {
		if gchannel.ParentID == parent.ID {
			s.ChannelDelete(gchannel.ID)
		}
	}

	msg = "Deleting category: " + teamName
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
	s.ChannelDelete(parent.ID)

	msg = "Deleting role: " + teamName
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
	for _, grole := range roles {
		if grole.ID == role.ID {
			s.GuildRoleDelete(guild.ID, grole.ID)
		}
	}
}

func findChannelByName(s *discordgo.Session, i *discordgo.InteractionCreate, channelName string) *discordgo.Channel {
	var channel *discordgo.Channel
	c, _ := s.Channel(i.ChannelID)
	guild, _ := s.Guild(c.GuildID)
	channels, _ := s.GuildChannels(guild.ID)
	for j := 0; j < len(channels); j++ {
		log.Printf("%s", (channels[j].Name))
		if channels[j].Name == channelName {
			channel = channels[j]
		}
	}
	return channel
}
