# Ubiquiti Device Monitor

A Go-based service for monitoring Ubiquiti network devices (routers, switches, cameras, door access systems) and providing REST API access to device status and diagnostics.

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
  - [Run with Docker](#run-with-docker)
  - [Run without Docker](#run-without-docker)
    - [Setting Up PostgreSQL](#setting-up-postgresql)
    - [Run the Application](#run-the-application)
- [API Endpoints](#api-endpoints)

## Features

- Device status monitoring
- REST API for device management
- PostgreSQL database for persistent storage
- Device status history tracking
- Support for various device types (routers, switches, cameras, door access systems)

## Getting Started

First, clone the repository:
```bash
git clone https://github.com/benjaminleon/ubiquiti-monitor.git
cd ubiquiti-monitor
```

### Run with Docker

The easiest way to run the application is using Docker Compose:

```bash
# Start the application and database
docker compose up --build
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

### Run without Docker
If you prefer to run the application with a local PostgreSQL database instead of using Docker, follow these steps (only verified on mac):

#### Setting Up PostgreSQL

1. Install PostgreSQL on your system:
   - **macOS**: `brew install postgresql`
   - **Ubuntu/Debian**: `sudo apt-get install postgresql`
   - **Windows**: Download from [PostgreSQL website](https://www.postgresql.org/download/windows/)

2. Start the PostgreSQL service:
   - **macOS**: `brew services start postgresql`
   - **Ubuntu/Debian**: `sudo service postgresql start`
   - **Windows**: The installer will set up the service automatically

3. Create the database and user:
```./setup_db.sh```

#### Run the Application

1. Install dependencies:
```bash
go mod download
```

2. Set up environment variables:
The environment variables need to be exported. Do something like:
```bash
set -a
source .env
set +a
```

3. Run the application:
```bash
go run main.go
```


## API Endpoints

### Devices

- `GET /api/devices` - List all devices, returning the latest status update for each unique device (identified by serial number). The `time_since_seen` field indicates how long it has been since each device's last status update, and will continue to grow if a device stops sending updates.
- `POST /api/devices` - Creates a new device info entry

Example POST request:
```bash
curl -X POST localhost:8080/api/devices/ \
  -H "Content-Type: application/json" \
  -d '{
    "ip_address": "192.168.1.3",
    "device_type": "camera",
    "serial_number": "1",
    "hardware_version": "1",
    "software_version": "1.12.22",
    "firmware_version": "1.0.0"
}'

curl -X POST localhost:8080/api/devices/ \
  -H "Content-Type: application/json" \
  -d '{
    "ip_address": "192.168.1.4",
    "device_type": "secret door",
    "serial_number": "2",
    "hardware_version": "0.1",
    "software_version": "1.2.3",
    "firmware_version": "4.5.6"
}'
```

Example GET request:
```bash
curl localhost:8080/api/devices/
[{"ip_address":"192.168.1.3","device_type":"muffin","serial_number":"1","hardware_version":"1","software_version":"1.12.22","firmware_version":"1.0.0","time_since_seen":"26m29s"},{"ip_address":"192.168.1.3","device_type":"banana","serial_number":"2","hardware_version":"1","software_version":"1.12.22","firmware_version":"1.0.0","time_since_seen":"2m37s"}]
```

