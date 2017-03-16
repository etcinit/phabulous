package requests

import "github.com/etcinit/gonduit/requests"

// PhabulousToSlackRequest represets a request to phabulous.toslack.
type PhabulousToSlackRequest struct {
	UserPHIDs []string `json:"UserPHIDs"`
	requests.Request
}
