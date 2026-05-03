package tlwpa4220

import "net/url"

// PowerLineStatistics holds the response from the powerline statistics endpoint.
type PowerLineStatistics struct {
	Timeout bool                      `json:"timeout"`
	Success bool                      `json:"success"`
	Data    []PowerLineStatisticsData `json:"data"`
}

// PowerLineStatisticsData holds data for a single paired powerline adapter.
type PowerLineStatisticsData struct {
	DeviceMac      string `json:"device_mac"`
	DevicePassword string `json:"device_password"`
	RxRate         string `json:"rx_rate"`
	TxRate         string `json:"tx_rate"`
	Status         string `json:"status"`
}

const (
	// PowerLineStatisticsPath is the endpoint for powerline statistics.
	PowerLineStatisticsPath string = "admin/powerline?form=plc_device"
)

var (
	// PowerLineStatisticsParams are the parameters to get powerline statistics.
	PowerLineStatisticsParams url.Values = url.Values{
		"operation": {"load"},
	}
)

// PowerLineStatistics returns the current powerline statistics from the device.
func (c Client) PowerLineStatistics() (PowerLineStatistics, error) {
	var powerLineStatistics PowerLineStatistics

	err := c.request(PowerLineStatisticsPath, PowerLineStatisticsParams, &powerLineStatistics)
	if err != nil {
		return PowerLineStatistics{}, err
	}

	return powerLineStatistics, nil
}
