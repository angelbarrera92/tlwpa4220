# TL-WPA4220 Prometheus Exporter

[![Lint Code Base](https://github.com/angelbarrera92/tlwpa4220/workflows/Lint%20Code%20Base/badge.svg)](https://github.com/angelbarrera92/tlwpa4220/actions/workflows/linter.yml)

Go library and Prometheus exporter for the TP-Link TL-WPA4220 powerline Wi-Fi extender.

## Device notes

The TL-WPA4220 runs a minimal embedded HTTP server with very limited concurrency. Sending rapid or concurrent requests will crash the web server, requiring a device reboot to recover. The exporter is designed around this constraint:

- Requests are made **sequentially**, never concurrently
- A configurable delay (`-request-delay`, default `1s`) is inserted between requests
- The scrape interval (`-interval`, default `30s`) is intentionally conservative

## Usage

### As a library

```go
package main

import (
    "fmt"

    "github.com/angelbarrera92/tlwpa4220/pkg/tlwpa4220"
)

func main() {
    t := tlwpa4220.Client{
        Username: "admin",
        Password: "yourpassword",
        IP:       "192.168.1.254",
    }

    wireless, err := t.WirelessStatistics()
    if err != nil {
        panic(err)
    }

    for _, device := range wireless.Data {
        fmt.Println(device.Mac, device.IP, device.Type)
    }

    powerline, err := t.PowerLineStatistics()
    if err != nil {
        panic(err)
    }

    for _, device := range powerline.Data {
        fmt.Println(device.DeviceMac, device.TxRate, device.RxRate)
    }
}
```

### As a CLI

```bash
# Build
make build
# or
go build -o tlwpa4220 cmd/main.go
```

Available subcommands:

```bash
# List connected wireless devices
./tlwpa4220 cli --device-ip=192.168.1.254 --password=yourpassword wireless-statistics

# List paired powerline adapters
./tlwpa4220 cli --device-ip=192.168.1.254 --password=yourpassword powerline-statistics

# Reboot the device
./tlwpa4220 cli --device-ip=192.168.1.254 --password=yourpassword reboot
```

### As a Prometheus exporter

```bash
./tlwpa4220 serve-metrics \
  --device-ip=192.168.1.254 \
  --password=yourpassword \
  --interval=30s \
  --request-delay=1s
```

All flags:

| Flag | Default | Description |
|---|---|---|
| `-device-ip` | *(required)* | IP address of the TL-WPA4220 |
| `-password` | *(required)* | Device admin password |
| `-username` | `admin` | Device admin username |
| `-host` | `0.0.0.0` | Host to listen on |
| `-port` | `8080` | Port to listen on |
| `-interval` | `30s` | Interval between full collection cycles |
| `-request-delay` | `1s` | Delay between consecutive HTTP requests to the device |

#### Endpoints

- `GET /metrics` — Prometheus metrics
- `GET /healthz` — Health check (returns `200 ok` when the exporter process is running)

#### Metrics

| Metric | Labels | Description |
|---|---|---|
| `scrape_success` | `endpoint` (`wireless`, `powerline`) | `1` if the last scrape succeeded, `0` otherwise |
| `connected_devices` | `mac`, `ip`, `devname`, `type`, `encryption` | `1` for each currently connected wireless device |
| `connected_devices_txpkts` | `mac`, `ip`, `devname`, `type`, `encryption` | Transmitted packets (always `0` on current firmware — known device bug) |
| `connected_devices_rxpkts` | `mac`, `ip`, `devname`, `type`, `encryption` | Received packets (always `0` on current firmware — known device bug) |
| `powerline_device` | `mac` | `1` if the powerline adapter is online, `0` otherwise |
| `powerline_device_tx_rate` | `mac` | Powerline TX rate in Mbps |
| `powerline_device_rx_rate` | `mac` | Powerline RX rate in Mbps |

> **Note:** `connected_devices_txpkts` and `connected_devices_rxpkts` are always `0` due to a firmware bug in the TL-WPA4220 — the device never populates these counters in its API response.

Example output:

```
# HELP scrape_success Whether the last scrape of the given endpoint succeeded (1=success, 0=failure)
# TYPE scrape_success gauge
scrape_success{endpoint="powerline"} 1
scrape_success{endpoint="wireless"} 1
# HELP connected_devices Wireless devices currently connected (1=connected)
# TYPE connected_devices gauge
connected_devices{devname="mylaptop",encryption="wpa2-psk",ip="192.168.1.100",mac="AA:BB:CC:DD:EE:FF",type="2.4GHz"} 1
# HELP powerline_device Powerline device status (1=on, 0=off)
# TYPE powerline_device gauge
powerline_device{mac="11:22:33:44:55:66"} 1
# HELP powerline_device_rx_rate Powerline device RX rate in Mbps
# TYPE powerline_device_rx_rate gauge
powerline_device_rx_rate{mac="11:22:33:44:55:66"} 254
# HELP powerline_device_tx_rate Powerline device TX rate in Mbps
# TYPE powerline_device_tx_rate gauge
powerline_device_tx_rate{mac="11:22:33:44:55:66"} 139
```

#### Useful PromQL queries

```promql
# Number of connected wireless devices
count(connected_devices)

# Alert when a scrape fails
scrape_success == 0

# Powerline TX rate for all adapters
powerline_device_tx_rate

# Connected devices by band
count(connected_devices) by (type)
```

### Docker

```bash
docker run --rm \
  -p 8080:8080 \
  ghcr.io/angelbarrera92/tlwpa4220:latest \
  serve-metrics \
  --device-ip=192.168.1.254 \
  --password=yourpassword
```

### Docker Compose (with Prometheus and Grafana)

A ready-to-use stack is provided in the `deploy/` directory:

```bash
cd deploy
# Edit docker-compose.yml to set your device IP and password
docker compose up -d
```

This starts:
- The exporter on port `8080`
- Prometheus on port `9090`
- Grafana on port `3000` (default credentials: `admin` / `admin`)

The Grafana dashboard is provisioned automatically and available under the **TL-WPA4220** folder.

## Development

Requirements: `make`, `docker`, `go`, `golangci-lint`

```bash
# Lint
make lint

# Build binary (via Docker)
make build

# Build container image
make container
```

## License

[GPLv3](LICENSE)
