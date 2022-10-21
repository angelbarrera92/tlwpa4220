# TL-WPA4220 go package

[![GitHub Super-Linter](https://github.com/angelbarrera92/tlwpa4220/workflows/Lint%20Code%20Base/badge.svg)](https://github.com/marketplace/actions/super-linter)

Go library for the TL-WPA4220

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
        Password: "hopeItsNotAdmin",
        IP:       "192.168.2.1",
    }

    wireless, err := t.WirelessStatistics()
    if err != nil {
        panic(err)
    }

    for _, device := range wireless.Data {
        fmt.Println(device.Mac)
    }

    err = t.Reboot()
    if err != nil {
        fmt.Println(err)
        panic(err)
    }
}
```

### As a CLI

This package also provides a CLI to interact with the TL-WPA4220.
It provides the following commands:

- `tlwpa4220 cli`: Interact with the TL-WPA4220 using a CLI
- `tlwpa4220 serve-metrics`: Expose metrics from the TL-WPA4220 in Prometheus format

#### Build

```bash
go build -o tlwpa4220 cmd/main.go
```

#### Run commands

The subcommands are:

- `wireless-statistics`: Get the wireless statistics
- `powerline-statistics`: Get the powerline statistics
- `reboot`: Reboot the device

```bash
$ ./tlwpa4220 cli --device-ip="192.168.2.1" --password="hopeItsNotAdmin" wireless-statistics
2022/10/21 17:48:42 Running command wireless-statistics
2022/10/21 17:48:42 Wireless statistics: {"success":true,"timeout":false,"data":[{"mac":"XX-YY-ZZ-46-1E-40","type":"2.4GHz","encryption":"wpa2-psk","rxpkts":"0","txpkts":"0","ip":"192.168.22.22","devName":"NSA"}],"others":{"max_rules":64}}
$ ./tlwpa4220 cli --device-ip="192.168.2.1" --password="hopeItsNotAdmin" reboot
$ ./tlwpa4220 cli --device-ip="192.168.2.1" --password="hopeItsNotAdmin" powerline-statistics
2022/10/21 21:23:57 Running command powerline-statistics
2022/10/21 21:23:57 Powerline statistics: {"timeout":false,"success":true,"data":[{"device_mac":"ZZ-YY-XX-2D-6C-40","device_password":"","rx_rate":"282","tx_rate":"186","status":"on"}]}

```

#### Run the metrics exporter

This package also provides a metrics exporter to expose the wireless statistics as Prometheus metrics.

```bash
$ ./tlwpa4220 serve-metrics --device-ip="192.168.2.1" --password="hopeItsNotAdmin"
2022/10/21 17:51:28 Starting the metric recording thread
2022/10/21 17:51:28 Serving metrics on 0.0.0.0:8080/metrics
```

##### Metrics

The metrics are:

- `connected_devices_total`: Number of connected devices
- `connected_devices`: Number of connected devices per device type
- `connected_devices_txpkts`: Number of transmitted packets per device
- `connected_devices_rxpkts`: Number of received packets per device
- `powerline_devices_total`: Number of powerline devices
- `powerline_device`: Powerline device information and status
- `powerline_device_txpkts`: Number of transmitted packets per powerline device
- `powerline_device_rxpkts`: Number of received packets per powerline device

```bash
$ curl http://localhost:8080/metrics
# HELP connected_devices Connected devices
# TYPE connected_devices gauge
connected_devices{devname="NSA",ip="192.168.22.22",mac="XX:YY:ZZ:46:1E:40"} 1
# HELP connected_devices_rxpkts The total number of received packets
# TYPE connected_devices_rxpkts gauge
connected_devices_rxpkts{devname="NSA",ip="192.168.22.22",mac="XX:YY:ZZ:46:1E:40"} 0
# HELP connected_devices_total The total number of connected devices
# TYPE connected_devices_total gauge
connected_devices_total 1
# HELP connected_devices_txpkts The total number of transmitted packets
# TYPE connected_devices_txpkts gauge
connected_devices_txpkts{devname="NSA",ip="192.168.22.22",mac="XX:YY:ZZ:46:1E:40"} 0
# HELP powerline_device The powerline device
# TYPE powerline_device gauge
powerline_device{mac="ZZ:YY:XX:2D:6C:40"} 1
# HELP powerline_device_rxpkts The powerline device rxpkts
# TYPE powerline_device_rxpkts gauge
powerline_device_rxpkts{mac="ZZ:YY:XX:2D:6C:40"} 283
# HELP powerline_device_txpkts The powerline device txpkts
# TYPE powerline_device_txpkts gauge
powerline_device_txpkts{mac="ZZ:YY:XX:2D:6C:40"} 178
# HELP powerline_devices_total The total number of powerline devices
# TYPE powerline_devices_total gauge
powerline_devices_total 1
```

## Development

Use the `Makefile` to build and lint the code.

Requirements:

- `make`
- [`Docker`](https://docs.docker.com/get-docker/)

```bash
make clean lint build
```

## License

[GPLv3](LICENSE)
