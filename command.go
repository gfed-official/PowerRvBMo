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
    }	

    commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
        "ping":   PingHandler,
        "revert": RevertHandler,
    }

    componentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
        "revert-vm": RevertVMSelectHandler,
    }
)
