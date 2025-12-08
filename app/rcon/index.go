package rcon

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/ethanburkett/goadmin/app/config"
	"github.com/ethanburkett/goadmin/app/parser"
	"github.com/ethanburkett/goadmin/app/rcon/commands"
)

type Client struct {
	Host     string
	Port     int
	Password string
	Timeout  time.Duration
	conn     *net.UDPConn
}

func NewClient(config *config.Config) *Client {
	return &Client{
		Host:     config.Server.Host,
		Port:     config.Server.Port,
		Password: config.Server.RconPassword,
		Timeout:  5 * time.Second,
	}
}

func (c *Client) Connect() error {
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", c.Host, c.Port))
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

func (c *Client) SendCommand(command string) (string, error) {
	if c.conn == nil {
		return "", fmt.Errorf("not connected")
	}

	packet := fmt.Sprintf("\xFF\xFF\xFF\xFFrcon %s %s", c.Password, command)

	_, err := c.conn.Write([]byte(packet))
	if err != nil {
		return "", err
	}

	_ = c.conn.SetReadDeadline(time.Now().Add(c.Timeout))

	buf := make([]byte, 4096)
	n, _, err := c.conn.ReadFromUDP(buf)
	if err != nil {
		return "", err
	}

	response := string(buf[:n])

	response = strings.TrimPrefix(response, "\xFF\xFF\xFF\xFFprint\n")

	return response, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

type StatusPlayer struct {
	ID           int    `json:"id"`
	Score        int    `json:"score"`
	Ping         int    `json:"ping"`
	Uuid         string `json:"uuid"`
	SteamID      string `json:"steamId"`
	Name         string `json:"name"`
	StrippedName string `json:"strippedName"`
	Address      string `json:"address"`
	QPort        int    `json:"qPort"`
	Rate         int    `json:"rate"`
}

type StatusResponse struct {
	Hostname string         `json:"hostname"`
	Version  string         `json:"version"`
	Address  string         `json:"address"`
	OS       string         `json:"os"`
	Type     string         `json:"type"`
	Map      string         `json:"map"`
	Players  []StatusPlayer `json:"players"`
}

func (c *Client) Status() (*StatusResponse, error) {
	response, err := c.SendCommand("status")
	if err != nil {
		return nil, err
	}

	status := &StatusResponse{
		Players: []StatusPlayer{},
	}

	lines := strings.Split(response, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "hostname:") {
			status.Hostname = strings.TrimSpace(strings.TrimPrefix(line, "hostname:"))
		} else if strings.HasPrefix(line, "version :") {
			status.Version = strings.TrimSpace(strings.TrimPrefix(line, "version :"))
		} else if strings.HasPrefix(line, "udp/ip  :") {
			status.Address = strings.TrimSpace(strings.TrimPrefix(line, "udp/ip  :"))
		} else if strings.HasPrefix(line, "os      :") {
			status.OS = strings.TrimSpace(strings.TrimPrefix(line, "os      :"))
		} else if strings.HasPrefix(line, "type    :") {
			status.Type = strings.TrimSpace(strings.TrimPrefix(line, "type    :"))
		} else if strings.HasPrefix(line, "map     :") {
			status.Map = strings.TrimSpace(strings.TrimPrefix(line, "map     :"))
		}
	}

	playerDataStarted := false
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "num score ping") {
			playerDataStarted = true
			continue
		}

		if strings.HasPrefix(line, "---") {
			continue
		}

		if playerDataStarted {
			player := parsePlayerLine(line)
			if player != nil {
				status.Players = append(status.Players, *player)
			}
		}
	}

	return status, nil
}

func parsePlayerLine(line string) *StatusPlayer {
	fields := strings.Fields(line)

	if len(fields) < 9 {
		return nil
	}

	player := &StatusPlayer{}

	if id, err := strconv.Atoi(fields[0]); err == nil {
		player.ID = id
	}

	if score, err := strconv.Atoi(fields[1]); err == nil {
		player.Score = score
	}

	if ping, err := strconv.Atoi(fields[2]); err == nil {
		player.Ping = ping
	}

	player.Uuid = fields[3]

	player.SteamID = fields[4]

	nameStart := 5
	addressIndex := -1

	for i := 5; i < len(fields); i++ {
		if strings.Contains(fields[i], ":") && strings.Contains(fields[i], ".") {
			addressIndex = i
			break
		}
	}
	if addressIndex > 6 {
		player.Name = strings.Join(fields[nameStart:addressIndex-1], " ")
	} else if addressIndex == 6 {
		player.Name = fields[5]
	}

	if addressIndex != -1 && addressIndex < len(fields) {
		// Strip port from IP address
		address := fields[addressIndex]
		if colonIdx := strings.LastIndex(address, ":"); colonIdx != -1 {
			player.Address = address[:colonIdx]
		} else {
			player.Address = address
		}

		if addressIndex+1 < len(fields) {
			if qport, err := strconv.Atoi(fields[addressIndex+1]); err == nil {
				player.QPort = qport
			}
		}

		if addressIndex+2 < len(fields) {
			if rate, err := strconv.Atoi(fields[addressIndex+2]); err == nil {
				player.Rate = rate
			}
		}
	}

	player.StrippedName = parser.StripColorCodes(player.Name)

	return player
}

func (c *Client) GetPlayer(playerID string) (*commands.DumpUserInfo, error) {
	response, err := c.SendCommand(fmt.Sprintf("dumpuser %s", playerID))
	if err != nil {
		return nil, err
	}

	return commands.ParseDumpUser(response)
}
