package rest

import (
	"io"
	"time"

	"github.com/ethanburkett/goadmin/app/config"
	"github.com/ethanburkett/goadmin/app/logger"
	"github.com/ethanburkett/goadmin/app/rcon"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Api struct {
	engine *gin.Engine
	config *config.Config
	rcon   *rcon.Client
}

type ApiResponse struct {
	Data     interface{} `json:"data"`
	Success  bool        `json:"success"`
	Messages []string    `json:"messages,omitempty"`
	Error    interface{} `json:"error,omitempty"`
}

func New(cfg *config.Config, rconClient *rcon.Client) *Api {
	gin.DefaultWriter = logger.GinWriter{}
	gin.DefaultErrorWriter = logger.GinWriter{}

	// gin.DisableConsoleColor()

	r := gin.New()

	isDev := cfg.Environment == "development"
	var origins []string
	if isDev {
		origins = []string{"http://localhost:5173"}
	} else {
		origins = []string{"http://localhost:5173"} // you can change this to your prod URL
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(gin.LoggerWithWriter(io.Writer(logger.GinWriter{})))
	r.Use(gin.Recovery())
	r.Use(ResponseMiddleware())

	api := &Api{
		engine: r,
		config: cfg,
		rcon:   rconClient,
	}

	r.GET("/health", func(c *gin.Context) {
		c.Set("data", gin.H{"status": "ok"})
		c.Status(200)
	})

	RegisterPlayerRoutes(r, api)
	RegisterStatusRoute(r, api)
	RegisterAuthRoutes(r, api)
	RegisterRconRoutes(r, api)
	RegisterRBACRoutes(r, api)
	RegisterGroupRoutes(r, api)
	RegisterCommandRoutes(r, api)
	RegisterReportRoutes(r, api)
	RegisterIamGodRoute(r, api)
	RegisterAuditRoutes(r, api)

	rconClient.SendCommand("say ^5GoAdmin ^7now serving.")

	return api
}

func (api *Api) Engine() *gin.Engine {
	return api.engine
}
