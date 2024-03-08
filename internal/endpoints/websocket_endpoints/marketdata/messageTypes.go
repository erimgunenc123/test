package marketdata

type marketDataRequest struct {
	Action     string         `json:"action"`
	Parameters map[string]any `json:"parameters"`
}
