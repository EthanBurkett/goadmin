package automessages

import (
	"fmt"
	"time"

	"github.com/ethanburkett/goadmin/app/plugins"
)

// AutoMessagesPlugin sends periodic messages to the server
type AutoMessagesPlugin struct {
	ctx      *plugins.PluginContext
	ticker   *time.Ticker
	stopChan chan bool
	messages []string
	interval time.Duration
}

// Metadata returns plugin information
func (p *AutoMessagesPlugin) Metadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		ID:           "auto-messages",
		Name:         "Auto Messages",
		Version:      "1.0.0",
		Author:       "GoAdmin Team",
		Description:  "Sends periodic messages to the server",
		Website:      "https://github.com/ethanburkett/goadmin",
		Dependencies: []string{},
		Permissions: []string{
			"rcon.execute",
		},
	}
}

// Init initializes the plugin
func (p *AutoMessagesPlugin) Init(ctx *plugins.PluginContext) error {
	p.ctx = ctx
	p.stopChan = make(chan bool)
	p.interval = 30 * time.Second
	p.messages = []string{
		"Welcome to the server!",
		"Join our Discord: discord.gg/example",
		"Report bugs with !report",
		"Check your stats with !stats",
	}
	return nil
}

// Start starts the plugin
func (p *AutoMessagesPlugin) Start() error {
	// Recreate stopChan in case the plugin was stopped and started again
	p.stopChan = make(chan bool)
	p.ticker = time.NewTicker(p.interval)
	messageIndex := 0

	// Register a command to view next message
	if p.ctx.CommandAPI != nil {
		p.ctx.CommandAPI.RegisterCommand(plugins.CommandDefinition{
			Name:        "nextmsg",
			Usage:       "nextmsg",
			Description: "Shows the next auto message",
			MinArgs:     0,
			MaxArgs:     0,
			Handler: func(playerName, playerGUID string, args []string) error {
				nextMessage := p.messages[messageIndex]
				if p.ctx.RCONAPI != nil {
					p.ctx.RCONAPI.SendCommand(fmt.Sprintf(`tell %s "^3Next message: ^7%s"`, playerName, nextMessage))
				}
				return nil
			},
		})
	}

	go func() {
		for {
			select {
			case <-p.ticker.C:
				if p.ctx.RCONAPI != nil {
					message := p.messages[messageIndex]
					p.ctx.RCONAPI.SendCommand(fmt.Sprintf(`say "^7%s"`, message))
					fmt.Printf("[AutoMessages] Sent: %s\n", message)
					messageIndex = (messageIndex + 1) % len(p.messages)
				}
			case <-p.stopChan:
				return
			}
		}
	}()

	fmt.Printf("[AutoMessages] Started (interval: %v)\n", p.interval)
	return nil
}

func (p *AutoMessagesPlugin) Stop() error {
	if p.ticker != nil {
		p.ticker.Stop()
	}

	// Unregister commands
	if p.ctx.CommandAPI != nil {
		p.ctx.CommandAPI.UnregisterCommand("nextmsg")
	}

	close(p.stopChan)
	fmt.Printf("[AutoMessages] Stopped\n")
	return nil
}

// Reload reloads the plugin configuration
func (p *AutoMessagesPlugin) Reload() error {
	fmt.Printf("[AutoMessages] Reloaded\n")
	return nil
}

// Register the plugin with the global registry
func init() {
	plugins.Registry.Register(&AutoMessagesPlugin{})
}
