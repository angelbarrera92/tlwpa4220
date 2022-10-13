package tlwpa4220

const (
	// Reboot Endpoint
	RebootEndpointPath string = "admin/reboot.json"
)

var (
	// RebootParams Parameters for reboot
	RebootParams map[string][]string = map[string][]string{
		"operation": {"write"},
	}
)

func (c Client) Reboot() error {
	err := c.request(RebootEndpointPath, RebootParams, nil)
	if err != nil {
		return err
	}
	return nil
}
