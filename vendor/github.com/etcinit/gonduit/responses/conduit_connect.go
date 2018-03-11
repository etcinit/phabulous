package responses

// ConduitConnectResponse represents the response from calling conduit.connect.
type ConduitConnectResponse struct {
	SessionKey   string `json:"sessionKey"`
	ConnectionID int64  `json:"connectionID"`
}
