package dev

import "github.com/etcinit/phabulous/app/interfaces"

// TestCommand allows one to send test messages to the feed channel.
type TestCommand struct{}

// GetUsage returns the usage of this command.
func (t *TestCommand) GetUsage() string {
	return "dev:feedTest"
}

// GetDescription returns the description of this command.
func (t *TestCommand) GetDescription() string {
	return "Prints a test message."
}

// GetMatchers returns the matchers for this command.
func (t *TestCommand) GetMatchers() []string {
	return []string{"dev:feedTest"}
}

// GetIMMatchers returns IM matchers for this command.
func (t *TestCommand) GetIMMatchers() []string {
	return []string{}
}

// GetMentionMatchers returns the channel mention matchers for this command.
func (t *TestCommand) GetMentionMatchers() []string {
	return []string{}
}

// GetHandler returns the handler for this command.
func (t *TestCommand) GetHandler() interfaces.Handler {
	return func(s interfaces.Bot, m interfaces.Message, matches []string) {
		s.PostOnFeed("This is a test message. Please ignore me.")
	}
}
