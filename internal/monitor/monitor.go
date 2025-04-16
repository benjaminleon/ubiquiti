package monitor

import (
	"fmt"
	"net"
	"time"

	"github.com/ben/ubiquiti-monitor/internal/database"
	"gorm.io/gorm"
)

type DeviceMonitor struct {
	db     *gorm.DB
	config struct {
		checkInterval time.Duration
		timeout       time.Duration
	}
}

func NewDeviceMonitor(db *gorm.DB) *DeviceMonitor {
	return &DeviceMonitor{
		db: db,
		config: struct {
			checkInterval time.Duration
			timeout       time.Duration
		}{
			checkInterval: 5 * time.Minute,
			timeout:       30 * time.Second,
		},
	}
}

func (m *DeviceMonitor) Start() {
	ticker := time.NewTicker(m.config.checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		m.checkDevices()
	}
}

func (m *DeviceMonitor) checkDevices() {
	var devices []database.Device
	if err := m.db.Find(&devices).Error; err != nil {
		fmt.Printf("Error fetching devices: %v\n", err)
		return
	}

	for _, device := range devices {
		status := m.checkDevice(device)
		m.updateDeviceStatus(device, status)
	}
}

func (m *DeviceMonitor) checkDevice(device database.Device) string {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:22", device.IPAddress), m.config.timeout)
	if err != nil {
		return "offline"
	}
	defer conn.Close()

	// In a real implementation, you would:
	// 1. Try to connect to the device
	// 2. Authenticate
	// 3. Query device information
	// 4. Update device details in the database

	return "online"
}

func (m *DeviceMonitor) updateDeviceStatus(device database.Device, status string) {
	now := time.Now()

	// Update device status
	device.Status = status
	device.LastSeen = now
	if err := m.db.Save(&device).Error; err != nil {
		fmt.Printf("Error updating device status: %v\n", err)
		return
	}

	// Create status history entry
	statusHistory := database.DeviceStatus{
		DeviceID:  device.ID,
		Status:    status,
		Timestamp: now,
		Details:   fmt.Sprintf("Device %s is %s", device.Hostname, status),
	}

	if err := m.db.Create(&statusHistory).Error; err != nil {
		fmt.Printf("Error creating status history: %v\n", err)
	}
}
