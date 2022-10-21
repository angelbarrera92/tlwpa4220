package tlwpa4220

import "net/url"

type PowerLineStatistics struct {
	Timeout bool                      `json:"timeout"`
	Success bool                      `json:"success"`
	Data    []powerLineStatisticsData `json:"data"`
}

type powerLineStatisticsData struct {
	DeviceMac      string `json:"device_mac"`
	DevicePassword string `json:"device_password"`
	RxRate         string `json:"rx_rate"`
	TxRate         string `json:"tx_rate"`
	Status         string `json:"status"`
}

const (
	// PowerLineStatisticsPath Endpoint
	PowerLineStatisticsPath string = "admin/powerline?form=plc_device"
)

var (
	// PowerLineStatisticsParams Parameters to get powerline statistics
	PowerLineStatisticsParams url.Values = url.Values{
		"operation": {"load"},
	}
)

//nolint:typecheck
func (c Client) PowerLineStatistics() (PowerLineStatistics, error) {
	var powerLineStatistics PowerLineStatistics

	err := c.request(PowerLineStatisticsPath, PowerLineStatisticsParams, &powerLineStatistics)
	if err != nil {
		return PowerLineStatistics{}, err
	}

	return powerLineStatistics, nil
}
