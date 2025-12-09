package models

import (
	"time"

	"github.com/ethanburkett/goadmin/app/database"
	"gorm.io/gorm"
)

// Server represents a game server instance
type Server struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"uniqueIndex;not null" json:"name"` // Friendly name (e.g., "Main Server", "EU Server")
	Host         string         `gorm:"not null" json:"host"`             // Server hostname/IP
	Port         int            `gorm:"not null" json:"port"`             // Server port
	RconPort     int            `gorm:"not null" json:"rconPort"`         // RCON port
	RconPassword string         `gorm:"not null" json:"-"`                // RCON password (excluded from JSON)
	GamesMpPath  string         `json:"gamesMpPath"`                      // Path to games_mp.log file
	IsActive     bool           `gorm:"default:true" json:"isActive"`     // Whether server is active
	IsDefault    bool           `gorm:"default:false" json:"isDefault"`   // Default server for operations
	Description  string         `json:"description"`                      // Server description
	Region       string         `json:"region"`                           // Server region (e.g., "US-East", "EU-West")
	MaxPlayers   int            `gorm:"default:0" json:"maxPlayers"`      // Max player count
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	// Relationships - link server to data
	TempBans       []TempBan        `gorm:"foreignKey:ServerID;constraint:OnDelete:SET NULL;" json:"tempBans,omitempty"`
	Reports        []Report         `gorm:"foreignKey:ServerID;constraint:OnDelete:SET NULL;" json:"reports,omitempty"`
	CommandHistory []CommandHistory `gorm:"foreignKey:ServerID;constraint:OnDelete:SET NULL;" json:"commandHistory,omitempty"`
	InGamePlayers  []InGamePlayer   `gorm:"foreignKey:ServerID;constraint:OnDelete:CASCADE;" json:"inGamePlayers,omitempty"`
	ServerStats    []ServerStats    `gorm:"foreignKey:ServerID;constraint:OnDelete:CASCADE;" json:"serverStats,omitempty"`
}

// CreateServer creates a new server instance
func CreateServer(name, host, rconPassword, gamesMpPath, description, region string, port, rconPort, maxPlayers int, isDefault bool) (*Server, error) {
	db := database.DB

	// If this is being set as default, unset any existing default
	if isDefault {
		db.Model(&Server{}).Where("is_default = ?", true).Update("is_default", false)
	}

	server := &Server{
		Name:         name,
		Host:         host,
		Port:         port,
		RconPort:     rconPort,
		RconPassword: rconPassword,
		GamesMpPath:  gamesMpPath,
		IsActive:     true,
		IsDefault:    isDefault,
		Description:  description,
		Region:       region,
		MaxPlayers:   maxPlayers,
	}

	err := db.Create(server).Error
	return server, err
}

// GetServerByID retrieves a server by ID
func GetServerByID(id uint) (*Server, error) {
	db := database.DB
	var server Server
	err := db.First(&server, id).Error
	return &server, err
}

// GetServerByName retrieves a server by name
func GetServerByName(name string) (*Server, error) {
	db := database.DB
	var server Server
	err := db.Where("name = ?", name).First(&server).Error
	return &server, err
}

// GetDefaultServer retrieves the default server
func GetDefaultServer() (*Server, error) {
	db := database.DB
	var server Server
	err := db.Where("is_default = ? AND is_active = ?", true, true).First(&server).Error
	return &server, err
}

// GetAllServers retrieves all servers
func GetAllServers() ([]Server, error) {
	db := database.DB
	var servers []Server
	err := db.Order("is_default DESC, name ASC").Find(&servers).Error
	return servers, err
}

// GetActiveServers retrieves all active servers
func GetActiveServers() ([]Server, error) {
	db := database.DB
	var servers []Server
	err := db.Where("is_active = ?", true).Order("is_default DESC, name ASC").Find(&servers).Error
	return servers, err
}

// UpdateServer updates a server's information
func UpdateServer(id uint, updates map[string]interface{}) error {
	db := database.DB

	// If setting this as default, unset other defaults first
	if isDefault, ok := updates["is_default"].(bool); ok && isDefault {
		db.Model(&Server{}).Where("id != ? AND is_default = ?", id, true).Update("is_default", false)
	}

	return db.Model(&Server{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteServer soft deletes a server
func DeleteServer(id uint) error {
	db := database.DB
	return db.Delete(&Server{}, id).Error
}

// SetAsDefault sets this server as the default
func (s *Server) SetAsDefault() error {
	db := database.DB
	return db.Transaction(func(tx *gorm.DB) error {
		// Unset all other defaults
		if err := tx.Model(&Server{}).Where("id != ?", s.ID).Update("is_default", false).Error; err != nil {
			return err
		}
		// Set this as default
		return tx.Model(s).Update("is_default", true).Error
	})
}

// Activate activates the server
func (s *Server) Activate() error {
	db := database.DB
	return db.Model(s).Update("is_active", true).Error
}

// Deactivate deactivates the server
func (s *Server) Deactivate() error {
	db := database.DB
	return db.Model(s).Update("is_active", false).Error
}

// TableName specifies the table name
func (Server) TableName() string {
	return "servers"
}
