package core

import (
	"fmt"

	"github.com/etcinit/gonduit/requests"
	"github.com/etcinit/phabulous/app/messages"
	"github.com/etcinit/phabulous/app/modules"
)

// MemeCommand allows users to create memes.
type MemeCommand struct{}

// GetUsage returns the usage of this command.
func (c *MemeCommand) GetUsage() string {
	return "meme <macro> \"upper\" \"lower\""
}

// GetDescription returns the description of this command.
func (c *MemeCommand) GetDescription() string {
	return "Generates a meme using the specified macro."
}

// GetMatchers returns the matchers for this command.
func (c *MemeCommand) GetMatchers() []string {
	return []string{
		"meme ([^ ]{1,128}) [\"“](.{1,128})[\"”] [\"“](.{1,128})[\"”]$",
	}
}

// GetIMMatchers returns IM matchers for this command.
func (c *MemeCommand) GetIMMatchers() []string {
	return []string{}
}

// GetMentionMatchers returns the channel mention matchers for this command.
func (c *MemeCommand) GetMentionMatchers() []string {
	return []string{}
}

// GetHandler returns the handler for this command.
func (c *MemeCommand) GetHandler() modules.Handler {
	return func(s modules.Service, m messages.Message, matches []string) {
		s.StartTyping(m.GetChannel())

		conn, err := s.GetGonduit()
		if err != nil {
			s.Excuse(m, err)
			return
		}

		res, err := conn.MacroCreateMeme(requests.MacroCreateMemeRequest{
			MacroName: matches[1],
			UpperText: matches[2],
			LowerText: matches[3],
		})
		if err != nil {
			s.Excuse(m, err)

			return
		}

		s.PostImage(
			m.GetChannel(),
			fmt.Sprintf("%s", res.URI),
			res.URI,
			messages.IconTasks,
			true,
		)
	}
}
