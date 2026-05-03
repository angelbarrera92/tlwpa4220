package tlwpa4220

import "net/url"

// WirelessStatistics holds the response from the wireless statistics endpoint.
type WirelessStatistics struct {
	Success bool                     `json:"success"`
	Timeout bool                     `json:"timeout"`
	Data    []WirelessStatisticsData `json:"data"`
	Others  wirelessStatisticsOthers `json:"others"`
}

// WirelessStatisticsData holds data for a single connected wireless device.
type WirelessStatisticsData struct {
	Mac        string `json:"mac"`
	Type       string `json:"type"`
	Encryption string `json:"encryption"`
	Rxpkts     string `json:"rxpkts"`
	Txpkts     string `json:"txpkts"`
	IP         string `json:"ip"`
	DevName    string `json:"devName"`
}

type wirelessStatisticsOthers struct {
	MaxRules int `json:"max_rules"`
}

const (
	// WirelessStatisticsPath is the endpoint for wireless statistics.
	WirelessStatisticsPath string = "admin/wireless?form=statistics"
)

var (
	// WirelessStatisticsParams are the parameters to get wireless statistics.
	WirelessStatisticsParams url.Values = url.Values{
		"operation": {"load"},
	}
)

// WirelessStatistics returns the current wireless statistics from the device.
func (c Client) WirelessStatistics() (WirelessStatistics, error) {
	var wirelessStatistics WirelessStatistics

	err := c.request(WirelessStatisticsPath, WirelessStatisticsParams, &wirelessStatistics)
	if err != nil {
		return WirelessStatistics{}, err
	}

	return wirelessStatistics, nil
}
