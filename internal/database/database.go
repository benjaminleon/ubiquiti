package database

import (
	"fmt"
	"time"

	"github.com/ben/ubiquiti-monitor/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Device struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	IPAddress       string    `gorm:"not null" json:"ip_address"`
	Hostname        string    `json:"hostname"`
	DeviceType      string    `gorm:"not null" json:"device_type"` // router, switch, camera, door_access
	Model           string    `json:"model"`
	SerialNumber    string    `json:"serial_number"`
	HardwareVersion string    `json:"hardware_version"`
	SoftwareVersion string    `json:"software_version"`
	FirmwareVersion string    `json:"firmware_version"`
	Status          string    `gorm:"not null" json:"status"` // online, offline, unknown
	LastSeen        time.Time `json:"last_seen"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type DeviceStatus struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	DeviceID  uint      `gorm:"not null" json:"device_id"`
	Status    string    `gorm:"not null" json:"status"`
	Timestamp time.Time `gorm:"not null" json:"timestamp"`
	Details   string    `gorm:"type:text" json:"details"`
}

func Init(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&Device{}, &DeviceStatus{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
