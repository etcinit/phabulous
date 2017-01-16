package modules

import (
	"regexp"

	"github.com/etcinit/phabulous/app/interfaces"
)

const (
	// RegularMatcherType is used for all regular messages (on channels
	// usually) that do not mention the bot.
	RegularMatcherType MatcherType = iota

	// IMMatcherType is used for all messages sent directly to the bot
	// (private messages).
	IMMatcherType

	// MentionMatcherType is used for messages in channels but that directly
	// mention the bot's username.
	MentionMatcherType
)

// MatcherType is the type of a matcher.
type MatcherType int

// HandlerTuple is a simple pairing of a regular expressions and a handler
// function to be called if it is matched with the input message.
type HandlerTuple struct {
	Pattern *regexp.Regexp
	Handler interfaces.Handler
}

// GetPattern returns the regular expression pattern in this tuple.
func (t HandlerTuple) GetPattern() *regexp.Regexp {
	return t.Pattern
}

// GetHandler returns the Handler in this tuple.
func (t HandlerTuple) GetHandler() interfaces.Handler {
	return t.Handler
}

// RegexBuilder is a function to be provided by Connector implementations while
// using CompileHeaders. It allows the Connector to modify any regular
// expressions before they are added to the matcher lists.
type RegexBuilder func(matcherType MatcherType, pattern string) *regexp.Regexp

// CompileHandlers is a utility function for Connector implementations to
// process all the modules/commands provided into a mapping of HandlerTuples.
// The first list returned is expressions to be matched on regular chats, and
// the second list has expressions to be matched on private/direct
// conversations with a bot.
func CompileHandlers(
	modules []interfaces.Module,
	builder RegexBuilder,
) ([]interfaces.HandlerTuple, []interfaces.HandlerTuple) {
	handlers := []interfaces.HandlerTuple{}
	imHandlers := []interfaces.HandlerTuple{}

	for _, module := range modules {
		for _, command := range module.GetCommands() {
			for _, rgx := range command.GetMatchers() {
				handlers = append(handlers, HandlerTuple{
					Pattern: builder(RegularMatcherType, rgx),
					Handler: command.GetHandler(),
				})
			}

			for _, rgx := range command.GetIMMatchers() {
				imHandlers = append(imHandlers, HandlerTuple{
					Pattern: builder(IMMatcherType, rgx),
					Handler: command.GetHandler(),
				})
			}

			for _, rgx := range command.GetMentionMatchers() {
				handlers = append(handlers, HandlerTuple{
					Pattern: builder(MentionMatcherType, rgx),
					Handler: command.GetHandler(),
				})
			}
		}
	}

	return handlers, imHandlers
}
