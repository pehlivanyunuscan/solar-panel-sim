# Solar Panel Simulator with Prometheus Metrics

A Go-based solar panel (MPPT) simulator that generates realistic daily solar power patterns and exposes them as Prometheus metrics. Includes pre-configured Grafana dashboards for visualization.

## 🌟 Features

- **Realistic Solar Patterns**: Generates daily solar power patterns based on sunrise/sunset simulation with natural variations
- **Multiple Sensor Metrics**: Exposes various MPPT sensor readings (voltage, current, temperature, SOC, etc.)
- **Multi-Panel Support**: Simulates multiple solar panels with independent patterns
- **Prometheus Integration**: Exports all metrics in Prometheus format
- **Grafana Dashboards**: Pre-configured dashboards for monitoring:
  - Daily Production & Consumption Summary
  - MPPT Solar Pole Detail Panel
  - Total Active Pole Count

## 📋 Prerequisites

- Docker & Docker Compose
- Go 1.24+ (for local development)

## 🚀 Quick Start

### Using Docker Compose (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd prom-custom-metric

# Start all services
docker-compose up -d
```

This will start:
- **solar-sim-1**: Solar simulator on port `8081`
- **solar-sim-2**: Solar simulator on port `8082`
- **Prometheus**: Metrics collection on port `9090`
- **Grafana**: Visualization on port `3000`

### Access the Services

| Service | URL | Credentials |
|---------|-----|-------------|
| Grafana | http://localhost:3000 | admin / admin |
| Prometheus | http://localhost:9090 | - |
| Metrics (sim-1) | http://localhost:8081/metrics | - |
| Metrics (sim-2) | http://localhost:8082/metrics | - |

## 📊 Metrics

The simulator exposes the following metrics under `mppt_values`:

### Sensor Metrics
| Sensor | Description |
|--------|-------------|
| `aku gerilimi` | Battery voltage |
| `panel gerilimi` | Panel voltage |
| `sarj akimi` | Charge current |
| `yuk akimi` | Load current |
| `sicaklik` | Temperature |
| `soc` | State of Charge |
| `yuk gucu` | Load power |
| `panel gucu` | Panel power |
| `yuk durum` | Load status |
| `aku tipi` | Battery type |
| `sarj durum` | Charge status |
| `kapi bilgisi` | Door information |
| `sarj gucu` | Charge power |
| `panel akim` | Panel current |

### Role Status Metrics
- `role 1` through `role 7` - Relay status indicators

## 🏗️ Project Structure

```
prom-custom-metric/
├── main.go                 # Main application entry point
├── docker-compose.yml      # Docker Compose configuration
├── Dockerfile              # Container build configuration
├── prometheus.yml          # Prometheus scrape configuration
├── go.mod                  # Go module dependencies
├── cmd/
│   └── check_pattern.go    # Pattern checking utility
├── grafana/
│   ├── dashboards/         # Pre-configured Grafana dashboards
│   └── provisioning/       # Grafana auto-provisioning configs
├── logging/
│   └── logging.go          # Logging utilities
├── models/
│   └── panel.go            # Panel data model
└── patterngen/
    └── pattern_generator.go # Daily pattern generation logic
```

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `LOG_LEVEL` | Logging level (DEBUG, INFO, WARN, ERROR) | `INFO` |
| `PANEL_SEED` | Seed for pattern generation (unique per container) | - |
| `GRAFANA_USER` | Grafana admin username | `admin` |
| `GRAFANA_PASSWORD` | Grafana admin password | `admin` |

### Panel Configuration

Panels are configured in `models/panel.go`:

```go
var PanelList = []Panel{
    {MaxPower: 3000, LastPatternDay: -1},
}
```

## 🌅 Solar Pattern Generation

The simulator generates realistic daily solar power patterns:

- **00:00 - 06:00**: No power (night)
- **06:00 - 08:00**: Slow morning rise
- **08:00 - 12:00**: Fast increase to peak
- **12:00 - 15:00**: Peak power with natural noise (~15% variation)
- **15:00 - 18:00**: Gradual decline
- **18:00 - 24:00**: No power (night)

Each day generates a new pattern with randomized parameters for realistic variation.

## 🔧 Local Development

```bash
# Install dependencies
go mod download

# Run the application
go run main.go

# The metrics endpoint will be available at http://localhost:8080/metrics
```

## 📝 API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /metrics` | Prometheus metrics endpoint |
| `GET /health` | Health check endpoint |

## 🐳 Docker Commands

```bash
# Build and start services
docker-compose up -d --build

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

## 📈 Grafana Dashboards

The project includes three pre-configured dashboards:

1. **Günlük Üretim & Tüketim Özeti** - Daily production and consumption summary
2. **MPPT Güneş Direği Detay Panosu** - Detailed MPPT solar pole panel
3. **Toplam Aktif Direk Sayısı** - Total active pole count

Dashboards are automatically provisioned when Grafana starts.
