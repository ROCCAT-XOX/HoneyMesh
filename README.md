# HoneyMesh

<p align="center">
  <img src="https://github.com/ROCCAT-XOX/HoneyMesh/blob/main/assets/profileExample.png?raw=true" alt="HoneyMesh Logo" width="200"/>
</p>

HoneyMesh is an advanced IoT system for beehive monitoring using ESP32-based mesh networking and real-time data analysis.

## Overview

HoneyMesh provides beekeepers with comprehensive monitoring capabilities by collecting and visualizing vital beehive data:

- Weight measurements tracking honey production
- Internal and external temperature readings
- Humidity and air quality monitoring
- WiFi signal strength analytics for mesh network optimization

## System Architecture

### Hardware Components

1. **Mesh Nodes**:
   - ESP32/ESP8266 microcontrollers with sensors attached to each beehive
   - Collect and transmit data: weight, temperature (internal/external), humidity, air quality
   - Configured with unique node IDs for identification

2. **Mesh Master**:
   - Dual ESP32 setup (RX/TX configuration)
   - One ESP32 handles mesh network communications (receiving data from nodes)
   - Second ESP32 operates as a web server, exposing data via JSON API
   - Includes OLED display for status monitoring

### Software Components

1. **Backend**:
   - Go-based web server using the Gin framework
   - MongoDB for data storage and querying
   - RESTful API endpoints for data access

2. **Frontend**:
   - Responsive dashboard built with Bootstrap and ECharts
   - Real-time data visualization with trend analysis
   - Time-period selection for historical data viewing

3. **Authentication**:
   - Secure user management system
   - Password hashing with bcrypt
   - Session-based authentication

## Features

- **Real-time Monitoring**: View current beehive metrics through an intuitive dashboard
- **Trend Analysis**: Track changes in weight, temperature, and other metrics over time
- **Multi-hive Support**: Monitor multiple beehives from a single interface
- **Weather Integration**: External weather data correlation with OpenWeatherMap API
- **Responsive Design**: Mobile-friendly interface for on-the-go monitoring

## Installation

### Prerequisites

- Go 1.19+
- MongoDB
- ESP32 development environment (Arduino IDE with ESP32 board support)
- Network with internet access (for weather data)

### Backend Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/ROCCAT-XOX/HoneyMesh.git
   cd HoneyMesh
   ```

2. Install Go dependencies:
   ```bash
   go mod download
   ```

3. Configure MongoDB:
   ```bash
   # Ensure MongoDB is running
   mongod --dbpath /path/to/data/db
   ```

4. Start the server:
   ```bash
   go run main.go
   ```
   The server will be available at http://localhost:8080

### Hardware Setup

1. Flash the appropriate firmware to your ESP32 devices:
   - Use `Arduino/ESP32-MeshNode.ino` for sensor nodes
   - Use `Arduino/ESP32-MeshMaster.ino` for the mesh master
   - Use `Arduino/ESP32-with-Webserver.ino` for the web server interface

2. Configure the WiFi credentials in each sketch
3. Deploy the nodes at your beehives

## Docker Deployment

HoneyMesh can be deployed using Docker for easier management:

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/ROCCAT-XOX/HoneyMesh.git
   cd HoneyMesh
   ```

2. Copy the example environment file:
   ```bash
   cp .env-example .env
   ```

3. Edit the `.env` file to configure your settings:
   ```
   # Set to true to generate sample data on startup
   GENERATE_SAMPLE_DATA=true
   ```

4. Build and start the containers:
   ```bash
   docker-compose up -d
   ```

5. Access the application at http://localhost:8080

### Configuration

#### Environment Variables

- `GENERATE_SAMPLE_DATA`: Set to `true` to generate sample data when the application starts.
- `MONGODB_URI`: MongoDB connection URI (defaults to `mongodb://mongodb:27017` in the compose file).

#### Volumes

The MongoDB data is stored in a named volume `mongodb-data` to persist your database between container restarts.

### Troubleshooting

#### Checking Logs

To view logs for the containers:

```bash
# View logs for the app container
docker-compose logs app

# View logs for the MongoDB container
docker-compose logs mongodb

# Follow logs in real-time
docker-compose logs -f
```

#### MongoDB Shell Access

To access the MongoDB shell:

```bash
docker exec -it honeymesh-mongodb mongosh
```

#### Restarting Services

To restart specific services:

```bash
docker-compose restart app
docker-compose restart mongodb
```

## Development

### Project Structure

- `main.go`: Application entry point
- `Router.go`: HTTP route definitions
- `analytics.go`: Data analysis functions
- `MongoDB.go`: Database connection and operations
- `Arduino/`: ESP32 firmware files
- `templates/`: HTML templates for web interface
- `Dockerfile`: Defines how the Go application container is built
- `docker-compose.yml`: Orchestrates the application and database containers
- `.env`: Contains environment variables for configuring the application

## License

[License information]

## Contributors

[Promedia GmbH]