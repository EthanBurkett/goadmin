package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ethanburkett/goadmin/app/commands"
	"github.com/ethanburkett/goadmin/app/config"
	"github.com/ethanburkett/goadmin/app/database"
	"github.com/ethanburkett/goadmin/app/logger"
	"github.com/ethanburkett/goadmin/app/models"
	"github.com/ethanburkett/goadmin/app/parser"
	"github.com/ethanburkett/goadmin/app/rcon"
	"github.com/ethanburkett/goadmin/app/rest"
	"github.com/ethanburkett/goadmin/app/watcher"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	logger.Init("backend", cfg.Environment == "development")
	defer logger.Log.Sync()

	logger.Info("GoAdmin starting...")

	database.Init()

	database.AutoMigrate(
		models.OfflinePlayer{},
		models.User{},
		models.Session{},
		models.Role{},
		models.Permission{},
		models.Setting{},
		models.CommandHistory{},
		models.ServerStats{},
		models.SystemStats{},
		models.PlayerStats{},
		models.Group{},
		models.InGamePlayer{},
		models.CustomCommand{},
		models.Report{},
		models.TempBan{},
		models.AuditLog{},
	)

	initializeSuperAdminRole()
	initializeViewerRole()
	initializeDefaultGroups()
	initializeDefaultCommands()

	rconClient := rcon.NewClient(cfg)
	err = rconClient.Connect()
	if err != nil {
		panic(err)
	}
	defer rconClient.Close()

	_, err = rconClient.Status()
	if err != nil {
		logger.Error("Error getting RCON status", zap.Error(err))
		return
	} else {
		logger.Info("RCON connection established successfully.")
	}

	restServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.RestPort),
		Handler: rest.New(cfg, rconClient).Engine(),
	}

	go func() {
		logger.Info(fmt.Sprintf("Starting API on 0.0.0.0:%d", cfg.RestPort))
		if err := restServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("REST API error", zap.Error(err))
		}
	}()

	go startGamesMpWatcher(cfg, rconClient)

	// Start stats collector
	statsCollector := watcher.NewStatsCollector(rconClient)
	statsCollector.Start()
	defer statsCollector.Stop()

	go startTempBanChecker()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := restServer.Shutdown(ctx); err != nil {
		logger.Error("API shutdown error", zap.Error(err))
	}

	logger.Info("GoAdmin stopped.")
}

func initializeSuperAdminRole() {
	superAdminRole, err := models.GetRoleByName("super_admin")
	if err != nil {
		superAdminRole, err = models.CreateRole("super_admin", "Super Administrator with all permissions")
		if err != nil {
			logger.Error("Failed to create super admin role", zap.Error(err))
			return
		}
		logger.Info("Created super_admin role")
	}

	permissions := []struct {
		name        string
		description string
	}{
		{"players.view", "View player information"},
		{"players.manage", "Manage players"},
		{"status.view", "View server status"},
		{"rcon.command", "Execute RCON commands"},
		{"rcon.kick", "Kick players"},
		{"rcon.ban", "Ban players"},
		{"rcon.say", "Send messages to server"},
		{"rcon.map", "Control map and game settings"},
		{"rcon.admin", "Server administration (exec, writeconfig, set)"},
		{"rbac.manage", "Manage roles and permissions"},
		{"users.delete", "Delete user accounts"},
		{"reports.view", "View player reports"},
		{"reports.action", "Take action on reports (ban, tempban, dismiss)"},
	}

	for _, perm := range permissions {
		permission, err := models.GetPermissionByName(perm.name)
		if err != nil {
			permission, err = models.CreatePermission(perm.name, perm.description)
			if err != nil {
				logger.Warn("Failed to create permission", zap.String("permission", perm.name), zap.Error(err))
				continue
			}
		}
		models.AddPermissionToRole(superAdminRole.ID, permission.ID)
	}

	logger.Info("Super admin role initialized with all permissions")
}

func initializeViewerRole() {
	viewerRole, err := models.GetRoleByName("viewer")
	if err != nil {
		viewerRole, err = models.CreateRole("viewer", "Basic viewer with read-only access")
		if err != nil {
			logger.Error("Failed to create viewer role", zap.Error(err))
			return
		}
		logger.Info("Created viewer role")
	}

	viewPermissions := []string{"players.view", "status.view"}
	for _, permName := range viewPermissions {
		permission, err := models.GetPermissionByName(permName)
		if err == nil {
			models.AddPermissionToRole(viewerRole.ID, permission.ID)
		}
	}

	logger.Info("Viewer role initialized with basic permissions")
}

func initializeDefaultGroups() {
	defaultGroups := []struct {
		name        string
		power       int
		permissions []string
		description string
	}{
		{
			name:        "Owner",
			power:       100,
			permissions: []string{"all"},
			description: "Server owner with full control",
		},
		{
			name:        "Admin",
			power:       50,
			permissions: []string{"kick", "ban", "map", "say"},
			description: "Server administrator",
		},
		{
			name:        "VIP",
			power:       10,
			permissions: []string{"say"},
			description: "VIP player",
		},
	}

	for _, g := range defaultGroups {
		groups, err := models.GetAllGroups()
		if err != nil {
			logger.Error("Failed to get groups", zap.Error(err))
			continue
		}

		exists := false
		for _, existingGroup := range groups {
			if existingGroup.Name == g.name {
				exists = true
				break
			}
		}

		if !exists {
			permissionsJSON, _ := json.Marshal(g.permissions)
			err := models.CreateGroup(g.name, g.power, string(permissionsJSON), g.description)
			if err != nil {
				logger.Warn("Failed to create default group", zap.String("group", g.name), zap.Error(err))
			} else {
				logger.Info(fmt.Sprintf("Created default group: %s (power: %d)", g.name, g.power))
			}
		}
	}
}

func initializeDefaultCommands() {
	defaultCommands := []struct {
		name        string
		usage       string
		description string
		rconCommand string
		minArgs     int
		maxArgs     int
		minPower    int
		permissions []string
		isBuiltIn   bool
	}{
		{
			name:        "kick",
			usage:       "!kick <player> [reason]",
			description: "Kick a player from the server",
			rconCommand: "clientkick {playerId:arg0} {argsFrom:1}",
			minArgs:     1,
			maxArgs:     -1,
			minPower:    50,
			permissions: []string{"kick"},
			isBuiltIn:   false,
		},
		{
			name:        "ban",
			usage:       "!ban <player> [reason]",
			description: "Ban a player from the server",
			rconCommand: "banClient {playerId:arg0} {argsFrom:1}",
			minArgs:     1,
			maxArgs:     -1,
			minPower:    50,
			permissions: []string{"ban"},
			isBuiltIn:   false,
		},
		{
			name:        "say",
			usage:       "!say <message>",
			description: "Send a message to all players",
			rconCommand: "say ^3[Admin] ^7{argsFrom:0}",
			minArgs:     1,
			maxArgs:     -1,
			minPower:    10,
			permissions: []string{"say"},
			isBuiltIn:   false,
		},
		{
			name:        "map",
			usage:       "!map <mapname>",
			description: "Change the current map",
			rconCommand: "map {arg0}",
			minArgs:     1,
			maxArgs:     1,
			minPower:    50,
			permissions: []string{"map"},
			isBuiltIn:   false,
		},
		{
			name:        "restart",
			usage:       "!restart",
			description: "Restart the current map",
			rconCommand: "map_restart",
			minArgs:     0,
			maxArgs:     0,
			minPower:    50,
			permissions: []string{"map"},
			isBuiltIn:   false,
		},
		{
			name:        "putgroup",
			usage:       "!putgroup <player> <group>",
			description: "Assign a player to a group (built-in Go function)",
			rconCommand: "",
			minArgs:     2,
			maxArgs:     2,
			minPower:    80,
			permissions: []string{"putgroup"},
			isBuiltIn:   true,
		},
		{
			name:        "groups",
			usage:       "!groups",
			description: "List all available groups (built-in Go function)",
			rconCommand: "",
			minArgs:     0,
			maxArgs:     0,
			minPower:    0,
			permissions: []string{},
			isBuiltIn:   true,
		},
		{
			name:        "mygroup",
			usage:       "!mygroup",
			description: "Show your current group and permissions (built-in Go function)",
			rconCommand: "",
			minArgs:     0,
			maxArgs:     0,
			minPower:    0,
			permissions: []string{},
			isBuiltIn:   true,
		},
		{
			name:        "adminlist",
			usage:       "!adminlist",
			description: "List all online admins (built-in Go function)",
			rconCommand: "",
			minArgs:     0,
			maxArgs:     0,
			minPower:    0,
			permissions: []string{},
			isBuiltIn:   true,
		},
		{
			name:        "help",
			usage:       "!help [page]",
			description: "Show available commands with pagination (built-in Go function)",
			rconCommand: "",
			minArgs:     0,
			maxArgs:     1,
			minPower:    0,
			permissions: []string{},
			isBuiltIn:   true,
		},
		{
			name:        "report",
			usage:       "!report <player> <reason>",
			description: "Report a player for admin review (built-in Go function)",
			rconCommand: "",
			minArgs:     2,
			maxArgs:     -1,
			minPower:    0,
			permissions: []string{},
			isBuiltIn:   true,
		},
		{
			name:        "tempban",
			usage:       "!tempban <player> <duration> <reason>",
			description: "Temporarily ban a player (built-in Go function). Duration: {number}{m/h/d/M/y} (e.g., 5m, 2h, 3d, 1M, 2y)",
			rconCommand: "",
			minArgs:     3,
			maxArgs:     -1,
			minPower:    80,
			permissions: []string{"tempban"},
			isBuiltIn:   true,
		},
	}

	for _, cmd := range defaultCommands {
		// Check if command already exists
		_, err := models.GetCustomCommand(cmd.name)
		if err != nil {
			// Command doesn't exist, create it
			permissionsJSON, _ := json.Marshal(cmd.permissions)
			err := models.CreateCustomCommand(
				cmd.name,
				cmd.usage,
				cmd.description,
				cmd.rconCommand,
				string(permissionsJSON),
				"both",
				cmd.minArgs,
				cmd.maxArgs,
				cmd.minPower,
				cmd.isBuiltIn,
			)
			if err != nil {
				logger.Warn("Failed to create default command", zap.String("command", cmd.name), zap.Error(err))
			} else {
				logger.Info(fmt.Sprintf("Created default command: !%s", cmd.name))
			}
		}
	}
}

func startGamesMpWatcher(cfg *config.Config, rconClient *rcon.Client) {
	changesChan := watcher.WatchGamesMp(cfg)

	cmdHandler := commands.NewCommandHandler(rconClient)

	for event := range changesChan {
		entry, ok := parser.ParseGamesMpLine(event.NewLine)
		if !ok {
			continue
		}

		switch entry.CommandType {
		case parser.SAY, parser.SAYTEAM:
			cleanMsg := strings.Map(func(r rune) rune {
				if r < 32 || r == 127 {
					return -1
				}
				return r
			}, entry.Message)
			cleanMsg = strings.TrimSpace(cleanMsg)

			if len(cleanMsg) > 0 && cleanMsg[0] == '!' {
				models.CreateOrUpdateInGamePlayer(entry.PlayerGUID, entry.PlayerName)

				if err := cmdHandler.ProcessChatCommand(entry.PlayerName, entry.PlayerGUID, cleanMsg); err != nil {
					logger.Error("Failed to process command", zap.Error(err))
				}
			}

		case parser.JOIN:
			fmt.Printf("[JOIN] %s (GUID: %s, ID: %s) joined the server\n", entry.PlayerName, entry.PlayerGUID, entry.PlayerID)
			models.CreateOrUpdateInGamePlayer(entry.PlayerGUID, entry.PlayerName)

			if models.IsPlayerTempBanned(entry.PlayerGUID) {
				ban, _ := models.GetTempBanByGUID(entry.PlayerGUID)
				if ban != nil {
					timeRemaining := time.Until(ban.ExpiresAt)
					hours := int(timeRemaining.Hours())
					minutes := int(timeRemaining.Minutes()) % 60
					kickMsg := fmt.Sprintf("You are temporarily banned. %dh %dm remaining. Reason: %s", hours, minutes, ban.Reason)
					rconClient.SendCommand(fmt.Sprintf("clientkick %s %s", entry.PlayerID, kickMsg))
					logger.Info(fmt.Sprintf("Kicked temp-banned player %s (%s)", entry.PlayerName, entry.PlayerGUID))
				}
			}

		case parser.LEAVE:
			fmt.Printf("[LEAVE] %s (GUID: %s, ID: %s) left the server\n", entry.PlayerName, entry.PlayerGUID, entry.PlayerID)
		}
	}
}

func startTempBanChecker() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := models.ExpireTempBans(); err != nil {
			logger.Error("Failed to expire temp bans", zap.Error(err))
		}
	}
}
