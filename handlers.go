package main

import (
    "fmt"
    "regexp"

	"github.com/bwmarrin/discordgo"
)

func RevertHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channel, _ := s.Channel(i.ChannelID)
	category, _ := s.Channel(channel.ParentID)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
            Flags: discordgo.MessageFlagsEphemeral,
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
                    CustomID: "revert-vm",
                    Placeholder: "Select a VM",
                    Options: buildSelectOptions(vms),
                },
            },
        },
    }

    contentMsg := fmt.Sprintf("Select a VM to revert to snapshot:")
    s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
        Content: &contentMsg,
        Components: &msgArray,
    })
}

func buildSelectOptions(vmList []string) []discordgo.SelectMenuOption {
    var options []discordgo.SelectMenuOption
	for _, vm := range vmList {
		options = append(options, discordgo.SelectMenuOption{
			Label: vm,
			Value: vm,
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

    success := fmt.Sprintf("VM reverted successfully!")

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
