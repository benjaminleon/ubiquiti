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
	DeviceType      string    `gorm:"not null" json:"device_type"` // router, switch, camera, door_access
	SerialNumber    string    `json:"serial_number"`
	HardwareVersion string    `json:"hardware_version"`
	SoftwareVersion string    `json:"software_version"`
	FirmwareVersion string    `json:"firmware_version"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func Init(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// AutoMigrate will create the necessary tables
	if err := db.AutoMigrate(&Device{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
