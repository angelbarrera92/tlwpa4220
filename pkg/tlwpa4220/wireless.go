package tlwpa4220

type WirelessStatistics struct {
	Success bool                     `json:"success"`
	Timeout bool                     `json:"timeout"`
	Data    []wirelessStatisticsData `json:"data"`
	Others  wirelessStatisticsOthers `json:"others"`
}

type wirelessStatisticsData struct {
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
	// WirelessStatisticsPath Endpoint
	WirelessStatisticsPath string = "admin/wireless?form=statistics"
)

var (
	// WirelessStatisticsParams Parameters to get wireless statistics
	WirelessStatisticsParams map[string][]string = map[string][]string{
		"operation": {"load"},
	}
)

func (c Client) WirelessStatistics() (WirelessStatistics, error) {
	var wirelessStatistics WirelessStatistics

	err := c.request(WirelessStatisticsPath, WirelessStatisticsParams, &wirelessStatistics)
	if err != nil {
		return WirelessStatistics{}, err
	}

	return wirelessStatistics, nil
}
