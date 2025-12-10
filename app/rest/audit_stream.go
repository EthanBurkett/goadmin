package rest

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/ethanburkett/goadmin/app/logger"
	"github.com/ethanburkett/goadmin/app/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from same origin and localhost
		return true // TODO: Implement proper origin checking in production
	},
}

// AuditStreamManager manages WebSocket connections for real-time audit log streaming
type AuditStreamManager struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan *models.AuditLog
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.RWMutex
}

// GlobalAuditStreamManager is the singleton instance
var GlobalAuditStreamManager *AuditStreamManager

// InitAuditStreamManager initializes the global audit stream manager
func InitAuditStreamManager() {
	GlobalAuditStreamManager = &AuditStreamManager{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan *models.AuditLog, 100), // Buffer for 100 messages
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
	go GlobalAuditStreamManager.run()
}

// run processes registration, unregistration, and broadcasting
func (asm *AuditStreamManager) run() {
	for {
		select {
		case conn := <-asm.register:
			asm.mu.Lock()
			asm.clients[conn] = true
			asm.mu.Unlock()
			logger.Info("Audit stream client connected",
				zap.String("remote_addr", conn.RemoteAddr().String()))

		case conn := <-asm.unregister:
			asm.mu.Lock()
			if _, ok := asm.clients[conn]; ok {
				delete(asm.clients, conn)
				conn.Close()
				logger.Info("Audit stream client disconnected",
					zap.String("remote_addr", conn.RemoteAddr().String()))
			}
			asm.mu.Unlock()

		case log := <-asm.broadcast:
			asm.mu.RLock()
			for conn := range asm.clients {
				// Send with timeout to prevent blocking
				go func(c *websocket.Conn, l *models.AuditLog) {
					if err := c.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
						logger.Error("Failed to set write deadline", zap.Error(err))
						asm.unregister <- c
						return
					}

					if err := c.WriteJSON(l); err != nil {
						logger.Error("Failed to write audit log to WebSocket",
							zap.Error(err),
							zap.String("remote_addr", c.RemoteAddr().String()))
						asm.unregister <- c
					}
				}(conn, log)
			}
			asm.mu.RUnlock()
		}
	}
}

// BroadcastAuditLog sends an audit log to all connected clients
func (asm *AuditStreamManager) BroadcastAuditLog(log *models.AuditLog) {
	select {
	case asm.broadcast <- log:
		// Successfully queued for broadcast
	default:
		// Channel is full, log a warning
		logger.Warn("Audit stream broadcast channel is full, dropping message",
			zap.Uint("log_id", log.ID))
	}
}

// handleAuditStream handles WebSocket connections for real-time audit log streaming
func handleAuditStream(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is authenticated (should have been set by AuthMiddleware)
		// The AuthMiddleware handles both Bearer tokens and session cookies
		_, exists := c.Get("user")
		if !exists {
			logger.Warn("Unauthenticated WebSocket connection attempt",
				zap.String("remote_addr", c.ClientIP()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}

		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logger.Error("Failed to upgrade WebSocket connection", zap.Error(err))
			c.Set("error", "Failed to upgrade connection")
			c.Status(http.StatusInternalServerError)
			return
		}

		logger.Info("WebSocket client connected",
			zap.String("remote_addr", conn.RemoteAddr().String()))

		// Register the connection
		GlobalAuditStreamManager.register <- conn

		// Handle client messages (mainly for keepalive)
		go func() {
			defer func() {
				GlobalAuditStreamManager.unregister <- conn
			}()

			// Set read deadline for pong messages
			conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			conn.SetPongHandler(func(string) error {
				conn.SetReadDeadline(time.Now().Add(60 * time.Second))
				return nil
			})

			for {
				// Read messages (we don't expect any, but this keeps the connection alive)
				var msg map[string]interface{}
				if err := conn.ReadJSON(&msg); err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						logger.Error("WebSocket read error", zap.Error(err))
					}
					break
				}

				// Handle ping messages
				if msgType, ok := msg["type"].(string); ok && msgType == "ping" {
					response := map[string]interface{}{
						"type":      "pong",
						"timestamp": time.Now(),
					}
					if err := conn.WriteJSON(response); err != nil {
						logger.Error("Failed to send pong", zap.Error(err))
						break
					}
				}
			}
		}()

		// Send ping messages to keep connection alive
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		go func() {
			for range ticker.C {
				if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
					logger.Error("Failed to send ping", zap.Error(err))
					GlobalAuditStreamManager.unregister <- conn
					return
				}
			}
		}()
	}
}

// StreamableAuditLog wraps the CreateAuditLog function to also broadcast to WebSocket clients
func StreamableAuditLog(log *models.AuditLog) {
	if GlobalAuditStreamManager != nil {
		GlobalAuditStreamManager.BroadcastAuditLog(log)
	}
}

// getAuditStreamStats returns statistics about connected WebSocket clients
func getAuditStreamStats(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		if GlobalAuditStreamManager == nil {
			c.Set("error", "Audit stream manager not initialized")
			c.Status(http.StatusServiceUnavailable)
			return
		}

		GlobalAuditStreamManager.mu.RLock()
		clientCount := len(GlobalAuditStreamManager.clients)
		GlobalAuditStreamManager.mu.RUnlock()

		stats := map[string]interface{}{
			"connected_clients": clientCount,
			"broadcast_buffer": map[string]interface{}{
				"capacity": cap(GlobalAuditStreamManager.broadcast),
				"size":     len(GlobalAuditStreamManager.broadcast),
			},
			"timestamp": time.Now(),
		}

		c.Set("data", stats)
		c.Status(http.StatusOK)
	}
}

// Helper function to convert AuditLog to JSON for streaming
func auditLogToJSON(log *models.AuditLog) ([]byte, error) {
	return json.Marshal(log)
}
