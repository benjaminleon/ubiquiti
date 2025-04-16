package api

import (
	"fmt"
	"net/http"

	"github.com/ben/ubiquiti-monitor/internal/config"
	"github.com/ben/ubiquiti-monitor/internal/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	router *gin.Engine
	db     *gorm.DB
	config config.ServerConfig
}

func NewServer(cfg config.ServerConfig, db *gorm.DB) *Server {
	router := gin.Default()
	server := &Server{
		router: router,
		db:     db,
		config: cfg,
	}

	// Setup routes
	server.setupRoutes()

	return server
}

func (s *Server) setupRoutes() {
	// Device endpoints
	devices := s.router.Group("/api/devices")
	{
		devices.GET("/", s.listDevices)
		devices.GET("/:id", s.getDevice)
		devices.GET("/:id/status", s.getDeviceStatus)
		devices.POST("/", s.createDevice)
		devices.PUT("/:id", s.updateDevice)
		devices.DELETE("/:id", s.deleteDevice)
	}
}

func (s *Server) Start() error {
	return s.router.Run(fmt.Sprintf(":%d", s.config.Port))
}

// Handler functions
func (s *Server) listDevices(c *gin.Context) {
	var devices []database.Device
	if err := s.db.Find(&devices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, devices)
}

func (s *Server) getDevice(c *gin.Context) {
	id := c.Param("id")
	var device database.Device
	if err := s.db.First(&device, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	c.JSON(http.StatusOK, device)
}

func (s *Server) getDeviceStatus(c *gin.Context) {
	id := c.Param("id")
	var status []database.DeviceStatus
	if err := s.db.Where("device_id = ?", id).Order("timestamp desc").Limit(10).Find(&status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, status)
}

func (s *Server) createDevice(c *gin.Context) {
	var device database.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.db.Create(&device).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, device)
}

func (s *Server) updateDevice(c *gin.Context) {
	id := c.Param("id")
	var device database.Device
	if err := s.db.First(&device, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.db.Save(&device).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, device)
}

func (s *Server) deleteDevice(c *gin.Context) {
	id := c.Param("id")
	if err := s.db.Delete(&database.Device{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
