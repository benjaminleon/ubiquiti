package api

import (
	"fmt"
	"net/http"
	"time"

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

// DeviceResponse represents the device data returned to the client
type DeviceResponse struct {
	IPAddress       string `json:"ip_address"`
	DeviceType      string `json:"device_type"`
	SerialNumber    string `json:"serial_number"`
	HardwareVersion string `json:"hardware_version"`
	SoftwareVersion string `json:"software_version"`
	FirmwareVersion string `json:"firmware_version"`
	TimeSinceSeen   string `json:"time_since_seen"`
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
		devices.POST("/", s.createDevice)
		devices.DELETE("/:id", s.deleteDevice)
	}
}

func (s *Server) Start() error {
	return s.router.Run(fmt.Sprintf(":%d", s.config.Port))
}

// Handler functions
func (s *Server) listDevices(c *gin.Context) {
	var devices []database.Device
	if err := s.db.
		Where("id IN (?)",
			s.db.Table("devices").
				Select("MAX(id)").
				Group("serial_number"),
		).
		Order("created_at DESC").
		Find(&devices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert devices to response format
	responses := make([]DeviceResponse, len(devices))
	for i, device := range devices {
		responses[i] = DeviceResponse{
			IPAddress:       device.IPAddress,
			DeviceType:      device.DeviceType,
			SerialNumber:    device.SerialNumber,
			HardwareVersion: device.HardwareVersion,
			SoftwareVersion: device.SoftwareVersion,
			FirmwareVersion: device.FirmwareVersion,
			TimeSinceSeen:   time.Since(device.CreatedAt).Round(time.Second).String(),
		}
	}

	c.JSON(http.StatusOK, responses)
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

func (s *Server) deleteDevice(c *gin.Context) {
	id := c.Param("id")
	if err := s.db.Delete(&database.Device{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
