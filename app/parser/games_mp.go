package parser

import (
	"regexp"
	"strings"
)

type CommandType int

const (
	SAY CommandType = iota
	SAYTEAM
	JOIN
	LEAVE
)

var Commands = map[string]CommandType{
	"say":     SAY,
	"sayteam": SAYTEAM,
	"J":       JOIN,
	"Q":       LEAVE,
}

type LogEntry struct {
	Command     string
	PlayerGUID  string
	PlayerID    string
	PlayerName  string
	Message     string
	CommandType CommandType
}

var timestampRegex = regexp.MustCompile(`^\s*\d+:\d+\s+`)

func ParseGamesMpLine(line string) (*LogEntry, bool) {
	if line == "" {
		return nil, false
	}

	line = timestampRegex.ReplaceAllString(line, "")
	line = strings.TrimSpace(line)

	if !strings.Contains(line, ";") {
		return nil, false
	}

	parts := strings.Split(line, ";")
	if len(parts) < 4 {
		return nil, false
	}

	command := parts[0]
	playerGUID := parts[1]
	playerID := parts[2]
	playerName := parts[3]

	cmdType, exists := Commands[command]
	if !exists {
		return nil, false
	}

	entry := &LogEntry{
		Command:     command,
		PlayerGUID:  playerGUID,
		PlayerID:    playerID,
		PlayerName:  playerName,
		CommandType: cmdType,
	}

	if cmdType == SAY || cmdType == SAYTEAM {
		if len(parts) >= 5 {
			entry.Message = parts[4]
		}
	}

	return entry, true
}
