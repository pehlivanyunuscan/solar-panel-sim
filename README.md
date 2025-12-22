# Solar Panel Simulator

A high-fidelity MPPT (Maximum Power Point Tracking) solar panel simulator written in Go. Generates realistic daily solar power generation patterns and exposes them as Prometheus metrics for monitoring and visualization.

## Overview

This simulator models the behavior of solar panels throughout a 24-hour cycle, producing realistic power output curves that account for sunrise, peak solar hours, and sunset. Each day generates a unique pattern with randomized parameters for natural variation, making it ideal for:

- Testing monitoring infrastructure
- Developing Grafana dashboards
- Learning Prometheus metrics collection
- Simulating IoT solar monitoring systems

## Features

- **Realistic Solar Patterns** — Mathematically modeled daily power curves with natural variations
- **Multi-Sensor Simulation** — Battery voltage, panel voltage, current, temperature, SOC, and more
- **Prometheus Native** — All metrics exposed in Prometheus format at `/metrics`
- **Pre-built Dashboards** — Ready-to-use Grafana dashboards included
- **Container-First** — Optimized for Docker/Kubernetes with stdout logging
- **Configurable** — Seed-based pattern generation for reproducible simulations

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.24+ (for local development)

### Run with Docker Compose

```bash
# Clone the repository
git clone https://github.com/yourusername/prom-custom-metric.git
cd prom-custom-metric

# Start all services
make docker-up

# Or using docker compose directly
docker compose up -d
```

### Access Services

| Service | URL | Credentials |
|---------|-----|-------------|
| Grafana | http://localhost:3000 | `admin` / `admin` |
| Prometheus | http://localhost:9090 | — |
| Simulator 1 | http://localhost:8081/metrics | — |
| Simulator 2 | http://localhost:8082/metrics | — |

## Architecture

```
┌─────────────────┐     ┌─────────────────┐
│  solar-sim-1    │     │  solar-sim-2    │
│  (Port 8081)    │     │  (Port 8082)    │
└────────┬────────┘     └────────┬────────┘
         │                       │
         └───────────┬───────────┘
                     │
              ┌──────▼──────┐
              │  Prometheus │
              │  (Port 9090)│
              └──────┬──────┘
                     │
              ┌──────▼──────┐
              │   Grafana   │
              │  (Port 3000)│
              └─────────────┘
```

## Project Structure

```
prom-custom-metric/
├── cmd/
│   ├── solar-sim/           # Main application
│   │   └── main.go
│   └── checkpattern/        # Pattern verification tool
│       └── main.go
├── internal/
│   ├── logging/             # Structured logging (stdout)
│   ├── models/              # Data models
│   └── patterngen/          # Solar pattern generator
├── configs/
│   └── prometheus.yml       # Prometheus configuration
├── grafana/
│   ├── dashboards/          # Pre-built dashboards
│   └── provisioning/        # Auto-provisioning configs
├── docs/
│   └── docker.md            # Docker documentation
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── go.mod
```

## Metrics

All metrics are exposed under the `mppt_values` gauge with a `sensor` label:

| Sensor | Description | Unit |
|--------|-------------|------|
| `panel gucu` | Panel power output (pattern-based) | Watts |
| `panel gerilimi` | Panel voltage | mV |
| `aku gerilimi` | Battery voltage | mV |
| `sarj akimi` | Charge current | mA |
| `yuk akimi` | Load current | mA |
| `sicaklik` | Temperature | °C |
| `soc` | State of Charge | % |
| `sarj gucu` | Charge power | Watts |
| `role 1-7` | Relay status | 0/1 |

### Example Prometheus Query

```promql
# Get current panel power
mppt_values{sensor="panel gucu"}

# Average power over last hour
avg_over_time(mppt_values{sensor="panel gucu"}[1h])
```

## Solar Pattern Generation

The simulator generates realistic daily solar power curves:

| Time | Phase | Behavior |
|------|-------|----------|
| 00:00 - 06:00 | Night | Zero output |
| 06:00 - 08:00 | Dawn | Gradual increase |
| 08:00 - 12:00 | Morning | Steep rise to peak |
| 12:00 - 15:00 | Peak | Maximum output with ±15% noise |
| 15:00 - 18:00 | Afternoon | Gradual decline |
| 18:00 - 24:00 | Night | Zero output |

Each day generates a new pattern with randomized:
- Peak hour timing (±30 min)
- Morning slope (18-22%)
- Peak noise (±7%)

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `LOG_LEVEL` | Logging level (`DEBUG`, `INFO`, `WARN`, `ERROR`) | `INFO` |
| `PANEL_SEED` | Seed for pattern generation (unique per instance) | `0` |
| `GRAFANA_USER` | Grafana admin username | `admin` |
| `GRAFANA_PASSWORD` | Grafana admin password | `admin` |

### Panel Configuration

Edit `internal/models/panel.go` to customize panel specifications:

```go
var PanelList = []Panel{
    {MaxPower: 3000, LastPatternDay: -1},
}
```

## Development

### Local Development

```bash
# Install dependencies
make deps

# Run the application
make run

# Run pattern verification
make test

# Build binary
make build

# Show all commands
make help
```

### Docker Commands

```bash
make docker-up       # Start all services
make docker-down     # Stop all services
make docker-logs     # View logs
make docker-ps       # Show container status
make docker-build    # Rebuild images
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/metrics` | GET | Prometheus metrics |
| `/health` | GET | Health check (returns JSON) |

### Health Check Response

```json
{
  "status": "healthy",
  "timestamp": "2025-01-15T10:30:00Z"
}
```

## Grafana Dashboards

Pre-configured dashboards are automatically provisioned:

1. **Daily Production & Consumption Summary** — Overview of energy metrics
2. **MPPT Solar Pole Detail Panel** — Detailed sensor readings
3. **Total Active Pole Count** — Fleet monitoring

## Documentation

- [Docker Guide](docs/docker.md) — Detailed Docker usage instructions