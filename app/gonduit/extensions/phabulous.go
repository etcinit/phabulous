package extensions

import (
	"github.com/etcinit/gonduit"
	"github.com/etcinit/phabulous/app/gonduit/extensions/requests"
	"github.com/etcinit/phabulous/app/gonduit/extensions/responses"
)

// PhabulousFromSlack performs a call to phabulous.fromslack.
func PhabulousFromSlack(
	c *gonduit.Conn,
	req requests.PhabulousFromSlackRequest,
) (*responses.PhabulousFromSlackResponse, error) {
	var res responses.PhabulousFromSlackResponse

	if err := c.Call("phabulous.fromslack", &req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// PhabulousToSlack performs a call to phabulous.toslack.
func PhabulousToSlack(
	c *gonduit.Conn,
	req requests.PhabulousToSlackRequest,
) (*responses.PhabulousToSlackResponse, error) {
	var res responses.PhabulousToSlackResponse

	if err := c.Call("phabulous.toslack", &req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
