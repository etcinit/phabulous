package bot

import (
	"fmt"

	"github.com/etcinit/gonduit/requests"
	"github.com/etcinit/phabulous/app/messages"
	"github.com/nlopes/slack"
)

// HandleCreateMeme create a meme and posts it to the channel.
func (b *Bot) HandleCreateMeme(ev *slack.MessageEvent, matches []string) {
	b.slackRTM.SendMessage(b.slackRTM.NewTypingMessage(ev.Channel))

	conn, err := b.Slacker.Factory.Make()
	if err != nil {
		b.Excuse(ev, err)
		return
	}

	res, err := conn.MacroCreateMeme(requests.MacroCreateMemeRequest{
		MacroName: matches[1],
		UpperText: matches[2],
		LowerText: matches[3],
	})
	if err != nil {
		b.Excuse(ev, err)

		return
	}

	b.Slacker.SimpleImagePost(
		ev.Channel,
		fmt.Sprintf("%s", res.URI),
		res.URI,
		messages.IconTasks,
		true,
	)
}
