package tlwpa4220

const (
	// RebootEndpointPath is the endpoint path for rebooting the device.
	RebootEndpointPath string = "admin/reboot.json"
)

var (
	// RebootParams are the parameters for the reboot request.
	RebootParams = map[string][]string{
		"operation": {"write"},
	}
)

// Reboot sends a reboot command to the device.
func (c Client) Reboot() error {
	return c.request(RebootEndpointPath, RebootParams, nil)
}
