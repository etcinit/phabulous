package bot

import (
	"github.com/etcinit/phabulous/app/messages"
	"github.com/nlopes/slack"
)

// HandleUsage shows usage tip.
func (b *Bot) HandleUsage(ev *slack.MessageEvent, matches []string) {
	b.Slacker.SimplePost(
		ev.Channel,
		"Hi. For usage information, type `help`.",
		messages.IconTasks,
	)
}

// HandleHelp shows help.
func (b *Bot) HandleHelp(ev *slack.MessageEvent, matches []string) {
	b.Slacker.SimplePost(
		ev.Channel,
		`Available commands:
    *lookup Txxx*: Looks up a task by its number.
    *lookup Dxxx*: Looks up a revision by its number.
    *help*: Shows this help.`,
		messages.IconTasks,
	)
}
