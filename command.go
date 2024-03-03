package main

import (
	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Replies with Pong!",
		},

		{
			Name:        "revert",
			Description: "Reverts a VM",
		},

        {
            Name:        "get",
            Description: "Get Reverts",
            Options: []*discordgo.ApplicationCommandOption{
                {
                    Type:        discordgo.ApplicationCommandOptionSubCommand,
                    Name:        "all",
                    Description: "Get all reverts",
                    Required:    false,
                },
                {
                    Type:        discordgo.ApplicationCommandOptionSubCommand,
                    Name:        "team",
                    Description: "Get reverts for a team",
                    Options: []*discordgo.ApplicationCommandOption{
                        {
                            Type:        discordgo.ApplicationCommandOptionString,
                            Name:        "team-id",
                            Description: "ID of the team",
                            Required:    true,
                        },
                    },
                },
            },
        },

		{
			Name:        "teams",
			Description: "Manage team pods",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "create",
					Description: "Subcommands group",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionInteger,
							Name:        "team-id",
							Description: "ID to start",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionInteger,
							Name:        "team-count",
							Description: "Number of teams to create",
							Required:    true,
						},
					},
				},
				{
					Name:        "delete",
					Description: "Subcommands group",
					Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "by-role",
							Description: "Specific team to delete",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
							Options: []*discordgo.ApplicationCommandOption{
								{
									Type:        discordgo.ApplicationCommandOptionRole,
									Name:        "team-role",
									Description: "Role for the team",
									Required:    true,
								},
							},
						},
						{
							Name:        "all",
							Description: "Delete all teams",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
						},
					},
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping":   PingHandler,
		"revert": RevertHandler,
		"teams":  TeamsHandler,
        "get":    GetRevertHandler,
	}

	componentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"revert-vm": RevertVMSelectHandler,
	}
)
