package interfaces

import "github.com/etcinit/phabulous/app/messages"

type Poster interface {
	Post(
		channelName string,
		storyText string,
		icon messages.Icon,
		asUser bool,
	)
	PostImage(
		channelName string,
		storyText string,
		imageURL string,
		icon messages.Icon,
		asUser bool,
	)
	PostOnFeed(storyText string) error
}
