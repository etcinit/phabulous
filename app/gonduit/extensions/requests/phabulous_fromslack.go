package requests

import "github.com/etcinit/gonduit/requests"

// PhabulousFromSlackRequest represents a requests to phabulous.fromslack.
type PhabulousFromSlackRequest struct {
	AccountIDs []string `json:"accountIDs"`
	requests.Request
}
