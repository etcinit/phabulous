package responses

import "github.com/etcinit/phabulous/app/gonduit/extensions/entities"

// PhabulousToSlackResponse is the response to calling phabulous.toslack.
type PhabulousToSlackResponse map[string]*entities.SlackUser
