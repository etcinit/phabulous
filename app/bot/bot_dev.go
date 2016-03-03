package bot

import "github.com/nlopes/slack"

func (b *Bot) HandleTestFeedMessage(ev *slack.MessageEvent, matches []string) {
	b.Slacker.FeedPost("This is a test")
}
