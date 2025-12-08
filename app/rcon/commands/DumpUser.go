package commands

import (
	"strconv"
	"strings"
)

type DumpUserInfo struct {
	IP             string `json:"ip"`
	PBGuid         string `json:"pbGuid"`
	XVer           string `json:"xVer"`
	QPort          int    `json:"qPort"`
	Challenge      int    `json:"challenge"`
	Protocol       int    `json:"protocol"`
	CgPredictItems int    `json:"cgPredictItems"`
	ClAnonymous    bool   `json:"clAnonymous"`
	ClPunkbuster   bool   `json:"clPunkbuster"`
	ClVoice        bool   `json:"clVoice"`
	ClWwwDownload  bool   `json:"clWwwDownload"`
	Rate           int    `json:"rate"`
	Snaps          int    `json:"snaps"`
	Name           string `json:"name"`
	PlayerID       string `json:"playerId"`
	PlayerSteamID  string `json:"playerSteamId"`
}

func ParseDumpUser(response string) (*DumpUserInfo, error) {
	info := &DumpUserInfo{}
	lines := strings.Split(response, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "userinfo" || strings.HasPrefix(line, "---") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		key := parts[0]
		value := strings.Join(parts[1:], " ")

		switch key {
		case "ip":
			if idx := strings.Index(value, "<=>"); idx != -1 {
				value = value[:idx]
			}
			// Strip port from IP address
			if colonIdx := strings.LastIndex(value, ":"); colonIdx != -1 {
				info.IP = value[:colonIdx]
			} else {
				info.IP = value
			}
		case "pbguid":
			info.PBGuid = value
		case "xver":
			info.XVer = value
		case "qport":
			if v, err := strconv.Atoi(value); err == nil {
				info.QPort = v
			}
		case "challenge":
			if v, err := strconv.Atoi(value); err == nil {
				info.Challenge = v
			}
		case "protocol":
			if v, err := strconv.Atoi(value); err == nil {
				info.Protocol = v
			}
		case "cg_predictItems":
			if v, err := strconv.Atoi(value); err == nil {
				info.CgPredictItems = v
			}
		case "cl_anonymous":
			info.ClAnonymous = value == "1"
		case "cl_punkbuster":
			info.ClPunkbuster = value == "1"
		case "cl_voice":
			info.ClVoice = value == "1"
		case "cl_wwwDownload":
			info.ClWwwDownload = value == "1"
		case "rate":
			if v, err := strconv.Atoi(value); err == nil {
				info.Rate = v
			}
		case "snaps":
			if v, err := strconv.Atoi(value); err == nil {
				info.Snaps = v
			}
		case "name":
			info.Name = value
		case "PlayerID":
			info.PlayerID = value
		case "PlayerSteamID":
			info.PlayerSteamID = value
		}
	}

	return info, nil
}
