package interfaces

import "github.com/etcinit/phabulous/app/messages"

// A Poster is an object capable of posting messages in a chat network.
type Poster interface {
	// Post posts a text message.
	Post(
		channelName string,
		storyText string,
		icon messages.Icon,
		asUser bool,
	)

	// PostImage posts a message with an attached image.
	PostImage(
		channelName string,
		storyText string,
		imageURL string,
		icon messages.Icon,
		asUser bool,
	)

	// PostOnFeed posts a message on the bot's "feed" channel.
	PostOnFeed(storyText string) error
}
