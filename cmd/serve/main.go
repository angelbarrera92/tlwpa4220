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

var (
	totalDevices = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "connected_devices_total",
		Help: "The total number of connected devices",
	})

	devices = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "connected_devices",
		Help: "Connected devices",
	}, []string{"mac", "ip", "devname"})

	txpkts = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "connected_devices_txpkts",
		Help: "The total number of transmitted packets",
	}, []string{"mac", "ip", "devname"})

	rxpkts = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "connected_devices_rxpkts",
		Help: "The total number of received packets",
	}, []string{"mac", "ip", "devname"})
)

type serveMetricsSubCommand struct {
	flag      *flag.FlagSet
	host      *string
	port      *int
	deviceIP  *string
	username  *string
	password  *string
	tlwpa4220 *tlwpa4220.Client
}

//nolint:revive
func NewServeMetricsSubCommand() *serveMetricsSubCommand {
	serveflag := flag.NewFlagSet("serve-metrics", flag.ExitOnError)

	host := serveflag.String("host", "0.0.0.0", "Host to serve metrics on")
	port := serveflag.Int("port", 8080, "Port to serve metrics on")
	deviceIP := serveflag.String("device-ip", "", "IP address of the TL-WPA4220 device to monitor")
	username := serveflag.String("username", "admin", "The username required to access the TL-WPA4220 device")
	password := serveflag.String("password", "", "The password required to access the TL-WPA4220 device")

	return &serveMetricsSubCommand{
		flag:     serveflag,
		host:     host,
		port:     port,
		deviceIP: deviceIP,
		username: username,
		password: password,
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

func (f *serveMetricsSubCommand) recordMetrics() {
	go func() {
		for {
			wirelessStatistics, err := f.tlwpa4220.WirelessStatistics()
			if err != nil {
				log.Printf("Error getting wireless statistics: %v", err)
			} else {
				totalDevices.Set(float64(len(wirelessStatistics.Data)))
				for _, data := range wirelessStatistics.Data {
					devices.WithLabelValues(formatMac(data.Mac), data.IP, data.DevName).Set(1)
					deviceTxpkts, err := strconv.ParseFloat(data.Txpkts, 64)
					if err != nil {
						log.Printf("Error parsing txpkts: %v", err)
					} else {
						txpkts.WithLabelValues(formatMac(data.Mac), data.IP, data.DevName).Set(deviceTxpkts)
					}
					deviceRxpkts, err := strconv.ParseFloat(data.Rxpkts, 64)
					if err != nil {
						log.Printf("Error parsing rxpkts: %v", err)
					} else {
						rxpkts.WithLabelValues(formatMac(data.Mac), data.IP, data.DevName).Set(deviceRxpkts)
					}
				}
			}
			// TODO Make this configurable
			time.Sleep(2 * time.Second)
		}
	}()
}

func formatMac(mac string) string {
	return strings.ReplaceAll(mac, "-", ":")
}

func (f *serveMetricsSubCommand) Run() error {
	// Run a simple web server to serve metrics
	http.Handle("/metrics", promhttp.Handler())
	listenAddress := fmt.Sprintf("%s:%d", *f.host, *f.port)
	log.Println("Starting the metric recording thread")
	f.recordMetrics()
	log.Printf("Serving metrics on %s/metrics", listenAddress)

	server := &http.Server{
		Addr:              listenAddress,
		ReadHeaderTimeout: 3 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
