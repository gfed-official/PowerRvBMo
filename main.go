package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
    "fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
)

var (
    GuildID       = flag.String("guild", "", "Guild ID")
    Token         = flag.String("token", "", "Bot Token")

	tomlConf      = Config{}
	vSphereClient *govmomi.Client
	ctx           = context.Background()
    finder        = &find.Finder{}

	RevertCounter = map[string]int{}
	RevertLimit   = 3
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
    var err error
    s, err = discordgo.New("Bot " + *Token)
    if err != nil {
        log.Fatalf("Cannot create a new session: %v", err)
    }

	ReadConfig(&tomlConf, "config.conf")
    fmt.Println("Connecting to vSphere...")
	vSphereClient = Connect()
	finder = find.NewFinder(vSphereClient.Client, true)
    dc, err := finder.Datacenter(ctx, tomlConf.VSphereDatacenter)
    finder.SetDatacenter(dc)
}

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
        switch i.Type {
        case discordgo.InteractionApplicationCommand:
            if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
                h(s, i)
            }
        case discordgo.InteractionMessageComponent:
            if h, ok := componentHandlers[i.MessageComponentData().CustomID]; ok {
                h(s, i)
            }
        }
	})
}

func main() {
    s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

    err := s.Open()
    if err != nil {
        log.Fatalf("Cannot open a new session: %v", err)
    }

    s.Identify.Intents = discordgo.IntentsGuildMessages

    fmt.Println(s.State)

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
		log.Printf("Added \"%s\"", v.Name)
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Gracefully shutting down.")
}
