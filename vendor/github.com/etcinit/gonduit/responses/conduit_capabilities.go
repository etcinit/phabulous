package responses

// ConduitCapabilitiesResponse represents a response from calling
// conduit.capabilities.
type ConduitCapabilitiesResponse struct {
	Authentication []string `json:"authentication"`
	Signatures     []string `json:"signatures"`
	Input          []string `json:"input"`
	Output         []string `json:"output"`
}
