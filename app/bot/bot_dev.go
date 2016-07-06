package bot

import "github.com/nlopes/slack"

// HandleTestFeedMessage is a handler that will post a simple test message to
// the feed channel.
func (b *Bot) HandleTestFeedMessage(ev *slack.MessageEvent, matches []string) {
	b.Slacker.FeedPost("This is a test")
}
