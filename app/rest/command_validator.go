package rest

import (
	"fmt"
	"regexp"
	"strings"
)

// CommandValidator validates and sanitizes RCON commands
type CommandValidator struct {
	// List of disallowed command prefixes
	disallowedCommands []string
	// List of blocked command patterns
	blockedPatterns []*regexp.Regexp
	// Maximum command length
	maxLength int
	// Maximum number of arguments
	maxArgs int
}

// NewCommandValidator creates a new command validator
func NewCommandValidator() *CommandValidator {
	return &CommandValidator{
		disallowedCommands: []string{
			// Server shutdown/control - dangerous
			"quit",
			"exit",
			"killserver",

			// Plugin/module loading - security risk
			"loadplugin",
			"unloadplugin",

			// Network changes - could break connectivity
			"net_restart",
			"net_ip",
			"net_port",

			// Developer/debug - should not be accessible
			"developer",
			"devmap",
			"sv_cheats",

			// Filesystem access - security risk
			"dir",
			"fs_game",
			"fs_homepath",
			"fs_basepath",

			// System commands
			"cmdlist",
			"cvarlist",
			"which",
			"vstr",
		},
		blockedPatterns: []*regexp.Regexp{
			// Block password exposure
			regexp.MustCompile(`(?i)rcon_password`),
			regexp.MustCompile(`(?i)sv_privatepassword`),
			regexp.MustCompile(`(?i)g_password`),
			// Block script execution that could be malicious
			regexp.MustCompile(`(?i);.*exec`), // Chained exec commands
			regexp.MustCompile(`(?i)\$\(`),    // Command substitution
			regexp.MustCompile(`(?i)&&`),      // Command chaining
			regexp.MustCompile(`(?i)\|\|`),    // OR chaining
			// Block potential injection attempts (pipe, redirect, backticks)
			regexp.MustCompile(`[|<>` + "`" + `]`), // Pipes, redirects, backticks only
		},
		maxLength: 500,
		maxArgs:   20,
	}
}

// ValidateCommand validates an RCON command
func (cv *CommandValidator) ValidateCommand(command string) error {
	// Check length
	if len(command) > cv.maxLength {
		return fmt.Errorf("command too long (max %d characters)", cv.maxLength)
	}

	// Check for empty command
	command = strings.TrimSpace(command)
	if command == "" {
		return fmt.Errorf("empty command")
	}

	// Check against blocked patterns
	for _, pattern := range cv.blockedPatterns {
		if pattern.MatchString(command) {
			return fmt.Errorf("command contains blocked pattern")
		}
	}

	// Split command into parts
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("invalid command format")
	}

	// Check argument count
	if len(parts)-1 > cv.maxArgs {
		return fmt.Errorf("too many arguments (max %d)", cv.maxArgs)
	}

	// Extract base command
	baseCmd := strings.ToLower(parts[0])

	// Check if command is in disallowed list
	for _, disallowedCmd := range cv.disallowedCommands {
		if baseCmd == disallowedCmd || strings.HasPrefix(baseCmd, disallowedCmd) {
			return fmt.Errorf("command not allowed: %s", baseCmd)
		}
	}

	return nil
}

// SanitizeCommand sanitizes a command string
func (cv *CommandValidator) SanitizeCommand(command string) string {
	// Remove null bytes
	command = strings.ReplaceAll(command, "\x00", "")
	// Remove carriage returns
	command = strings.ReplaceAll(command, "\r", "")
	// Replace newlines with spaces
	command = strings.ReplaceAll(command, "\n", " ")
	// Trim whitespace
	command = strings.TrimSpace(command)
	// Collapse multiple spaces
	command = regexp.MustCompile(`\s+`).ReplaceAllString(command, " ")

	return command
}

// IsRestrictedCommand checks if a command requires elevated permissions
func (cv *CommandValidator) IsRestrictedCommand(command string) bool {
	restrictedCommands := []string{
		"exec",
		"writeconfig",
		"set",
		"seta",
		"sets",
		"setu",
		"quit",
		"killserver",
		"map_restart",
		"fast_restart",
	}

	parts := strings.Fields(command)
	if len(parts) == 0 {
		return false
	}

	baseCmd := strings.ToLower(parts[0])
	for _, restricted := range restrictedCommands {
		if baseCmd == restricted {
			return true
		}
	}

	return false
}

// ValidateAndSanitize performs both validation and sanitization
func (cv *CommandValidator) ValidateAndSanitize(command string) (string, error) {
	// Sanitize first
	sanitized := cv.SanitizeCommand(command)

	// Then validate
	if err := cv.ValidateCommand(sanitized); err != nil {
		return "", err
	}

	return sanitized, nil
}

// Global validator instance
var CommandValidatorInstance = NewCommandValidator()

// ValidateRconCommand is a helper function to validate RCON commands
func ValidateRconCommand(command string) (string, error) {
	return CommandValidatorInstance.ValidateAndSanitize(command)
}
