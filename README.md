# Ubiquiti Device Monitor

A Go-based service for monitoring Ubiquiti network devices (routers, switches, cameras, door access systems) and providing REST API access to device status and diagnostics.

## Features

- Device status monitoring
- REST API for device management
- PostgreSQL database for persistent storage
- Device status history tracking
- Support for various device types (routers, switches, cameras, door access systems)

## Prerequisites

- Go 1.24.2 or later (for local development)
- Docker and Docker Compose (for containerized deployment)
- Network access to monitored devices

## Quick Start with Docker

The easiest way to run the application is using Docker Compose:

```bash
# Clone the repository
git clone https://github.com/ben/ubiquiti-monitor.git
cd ubiquiti-monitor

# Start the application and database
docker-compose up -d
```

This will:
- Start a PostgreSQL database container
- Build and start the application container
- Set up all necessary environment variables
- Create a persistent volume for the database

The application will be available at `http://localhost:8080`

To stop the application:
```bash
docker-compose down
```

## Local Development

1. Clone the repository:
```bash
git clone https://github.com/ben/ubiquiti-monitor.git
cd ubiquiti-monitor
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=ubiquiti_monitor
```

4. Run the application:
```bash
go run main.go
```

## API Endpoints

### Devices

- `GET /api/devices` - List all devices
- `GET /api/devices/:id` - Get device details
- `GET /api/devices/:id/status` - Get device status history
- `POST /api/devices` - Create a new device
- `PUT /api/devices/:id` - Update device details
- `DELETE /api/devices/:id` - Delete a device

## Database Schema

### Device
- ID (Primary Key)
- IP Address
- Hostname
- Device Type
- Model
- Serial Number
- Hardware Version
- Software Version
- Firmware Version
- Status
- Last Seen
- Created At
- Updated At

### Device Status
- ID (Primary Key)
- Device ID (Foreign Key)
- Status
- Timestamp
- Details

## Configuration

The application can be configured through environment variables:

- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database user (default: postgres)
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name (default: ubiquiti_monitor)
