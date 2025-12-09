package plugins

import (
	"fmt"
	"time"

	"github.com/ethanburkett/goadmin/app/rcon"
)

// RCONAPIImpl implements the RCONAPI interface for plugins
type RCONAPIImpl struct {
	client *rcon.Client
}

// NewRCONAPI creates a new RCON API instance
func NewRCONAPI(client *rcon.Client) *RCONAPIImpl {
	return &RCONAPIImpl{
		client: client,
	}
}

// SendCommand sends a raw RCON command
func (r *RCONAPIImpl) SendCommand(command string) (string, error) {
	if r.client == nil {
		return "", fmt.Errorf("RCON client not initialized")
	}
	return r.client.SendCommand(command)
}

// SendCommandWithTimeout sends a command with a custom timeout
func (r *RCONAPIImpl) SendCommandWithTimeout(command string, timeout time.Duration) (string, error) {
	if r.client == nil {
		return "", fmt.Errorf("RCON client not initialized")
	}
	return r.client.SendCommandWithTimeout(command, timeout)
}

// GetStatus gets server status
func (r *RCONAPIImpl) GetStatus() (map[string]interface{}, error) {
	if r.client == nil {
		return nil, fmt.Errorf("RCON client not initialized")
	}

	response, err := r.client.SendCommand("status")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"raw": response,
	}, nil
}
