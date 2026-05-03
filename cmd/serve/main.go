// Package serve provides the serve-metrics subcommand for exposing TL-WPA4220 metrics via Prometheus.
package serve

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/angelbarrera92/tlwpa4220/pkg/tlwpa4220"
)

const (
	defaultPort         = 8080
	defaultInterval     = 30 * time.Second
	defaultRequestDelay = time.Second
	readHeaderTimeout   = 3 * time.Second
)

var (
	powerLineDevices = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "powerline_device",
		Help: "Powerline device status (1=on, 0=off)",
	}, []string{"mac"})

	powerLineDeviceTxpkts = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "powerline_device_tx_rate",
		Help: "Powerline device TX rate in Mbps",
	}, []string{"mac"})

	powerLineDeviceRxpkts = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "powerline_device_rx_rate",
		Help: "Powerline device RX rate in Mbps",
	}, []string{"mac"})

	devices = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "connected_devices",
		Help: "Wireless devices currently connected (1=connected)",
	}, []string{"mac", "ip", "devname", "type", "encryption"})

	txpkts = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "connected_devices_txpkts",
		Help: "Transmitted packets per connected wireless device",
	}, []string{"mac", "ip", "devname", "type", "encryption"})

	rxpkts = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "connected_devices_rxpkts",
		Help: "Received packets per connected wireless device",
	}, []string{"mac", "ip", "devname", "type", "encryption"})

	scrapeSuccess = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "scrape_success",
		Help: "Whether the last scrape of the given endpoint succeeded (1=success, 0=failure)",
	}, []string{"endpoint"})
)

type serveMetricsSubCommand struct {
	flag         *flag.FlagSet
	host         *string
	port         *int
	deviceIP     *string
	username     *string
	password     *string
	interval     *time.Duration
	requestDelay *time.Duration
	tlwpa4220    *tlwpa4220.Client
}

//nolint:revive
func NewServeMetricsSubCommand() *serveMetricsSubCommand {
	serveflag := flag.NewFlagSet("serve-metrics", flag.ExitOnError)

	host := serveflag.String("host", "0.0.0.0", "Host to serve metrics on")
	port := serveflag.Int("port", defaultPort, "Port to serve metrics on")
	deviceIP := serveflag.String("device-ip", "", "IP address of the TL-WPA4220 device to monitor")
	username := serveflag.String("username", "admin", "The username required to access the TL-WPA4220 device")
	password := serveflag.String("password", "", "The password required to access the TL-WPA4220 device")
	interval := serveflag.Duration("interval", defaultInterval, "Interval between metric collection cycles")
	requestDelay := serveflag.Duration(
		"request-delay",
		defaultRequestDelay,
		"Delay between consecutive HTTP requests to the device (prevents overloading the embedded web server)",
	)

	return &serveMetricsSubCommand{
		flag:         serveflag,
		host:         host,
		port:         port,
		deviceIP:     deviceIP,
		username:     username,
		password:     password,
		interval:     interval,
		requestDelay: requestDelay,
	}
}

func (f *serveMetricsSubCommand) Parse(args []string) error {
	err := f.flag.Parse(args)
	if err != nil {
		return err
	}

	if *f.deviceIP == "" {
		return fmt.Errorf("device-ip is required")
	}

	if *f.password == "" {
		return fmt.Errorf("password is required")
	}

	f.tlwpa4220 = &tlwpa4220.Client{
		IP:       *f.deviceIP,
		Username: *f.username,
		Password: *f.password,
	}

	return nil
}

// Run starts the metrics HTTP server and the background collection loop.
func (f *serveMetricsSubCommand) Run() error {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	listenAddress := fmt.Sprintf("%s:%d", *f.host, *f.port)
	log.Println("Starting the metric recording thread")
	f.recordMetrics()
	log.Printf("Serving metrics on %s/metrics", listenAddress)
	log.Printf("Health endpoint on %s/healthz", listenAddress)

	server := &http.Server{
		Addr:              listenAddress,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (f *serveMetricsSubCommand) recordMetrics() {
	go func() {
		for {
			powerLineStatistics, err := f.tlwpa4220.PowerLineStatistics()
			if err != nil {
				log.Printf("Error getting powerline statistics: %s", err)
				scrapeSuccess.WithLabelValues("powerline").Set(0)
			} else {
				powerLineDevices.Reset()
				powerLineDeviceTxpkts.Reset()
				powerLineDeviceRxpkts.Reset()
				for _, device := range powerLineStatistics.Data {
					if device.Status == "on" {
						powerLineDevices.WithLabelValues(formatMac(device.DeviceMac)).Set(1)
					} else {
						powerLineDevices.WithLabelValues(formatMac(device.DeviceMac)).Set(0)
					}
					pLDeviceTxpkts, err := strconv.ParseFloat(device.TxRate, 64)
					if err != nil {
						log.Printf("Error parsing powerline tx_rate: %s", err)
					} else {
						powerLineDeviceTxpkts.WithLabelValues(formatMac(device.DeviceMac)).Set(pLDeviceTxpkts)
					}
					plDeviceRxpkts, err := strconv.ParseFloat(device.RxRate, 64)
					if err != nil {
						log.Printf("Error parsing powerline rx_rate: %s", err)
					} else {
						powerLineDeviceRxpkts.WithLabelValues(formatMac(device.DeviceMac)).Set(plDeviceRxpkts)
					}
				}
				scrapeSuccess.WithLabelValues("powerline").Set(1)
			}

			// Delay between requests to avoid overloading the device's embedded web server
			time.Sleep(*f.requestDelay)

			wirelessStatistics, err := f.tlwpa4220.WirelessStatistics()
			if err != nil {
				log.Printf("Error getting wireless statistics: %v", err)
				scrapeSuccess.WithLabelValues("wireless").Set(0)
			} else {
				devices.Reset()
				txpkts.Reset()
				rxpkts.Reset()
				for _, data := range wirelessStatistics.Data {
					devices.WithLabelValues(formatMac(data.Mac), data.IP, data.DevName, data.Type, data.Encryption).Set(1)
					deviceTxpkts, err := strconv.ParseFloat(data.Txpkts, 64)
					if err != nil {
						log.Printf("Error parsing txpkts: %v", err)
					} else {
						txpkts.WithLabelValues(formatMac(data.Mac), data.IP, data.DevName, data.Type, data.Encryption).Set(deviceTxpkts)
					}
					deviceRxpkts, err := strconv.ParseFloat(data.Rxpkts, 64)
					if err != nil {
						log.Printf("Error parsing rxpkts: %v", err)
					} else {
						rxpkts.WithLabelValues(formatMac(data.Mac), data.IP, data.DevName, data.Type, data.Encryption).Set(deviceRxpkts)
					}
				}
				scrapeSuccess.WithLabelValues("wireless").Set(1)
			}
			time.Sleep(*f.interval)
		}
	}()
}

func formatMac(mac string) string {
	return strings.ReplaceAll(mac, "-", ":")
}
