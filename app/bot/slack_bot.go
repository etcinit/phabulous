package bot

import (
	"math/rand"
	"regexp"

	"github.com/etcinit/gonduit"
	"github.com/etcinit/phabulous/app/factories"
	"github.com/etcinit/phabulous/app/interfaces"
	"github.com/etcinit/phabulous/app/messages"
	"github.com/etcinit/phabulous/app/modules/core"
	"github.com/etcinit/phabulous/app/modules/dev"
	"github.com/etcinit/phabulous/app/modules/extension"
	"github.com/jacobstr/confer"
	"github.com/nlopes/slack"
)

// NewBot creates a new instance of a Bot.
func NewSlackBot(
	slacker *SlackService,
	slackRTM *slack.RTM,
	slackInfo *slack.Info,
) *SlackBot {
	bot := &SlackBot{
		Slacker:      slacker,
		slackInfo:    slackInfo,
		slackRTM:     slackRTM,
		imChannelIDs: map[string]bool{},
		modules: []interfaces.Module{
			&dev.Module{},
			&core.Module{},
			&extension.Module{},
		},
	}

	// Make it easy to lookup if a channel is an IM channel.
	for _, im := range slackInfo.IMs {
		bot.imChannelIDs[im.ID] = true
	}

	// Load message handlers
	bot.loadHandlers()

	return bot
}

// SlackBot represents the state of the bot. It also contains functions
// related to the interactive portion of Phabulous.
type SlackBot struct {
	Slacker *SlackService
	Factory *factories.GonduitFactory

	slackInfo    *slack.Info
	slackRTM     *slack.RTM
	imChannelIDs map[string]bool
	handlers     []HandlerTuple
	imHandlers   []HandlerTuple

	modules []interfaces.Module
}

// HandlerTuple a tuples of a pattern and a handler.
type HandlerTuple struct {
	Pattern *regexp.Regexp
	Handler interfaces.Handler
}

func (b *SlackBot) mentionRegex(contents string) *regexp.Regexp {
	username := b.slackInfo.User.ID

	return regexp.MustCompile("^<@" + username + ">:? " + contents + "$")
}

func (b *SlackBot) loadHandlers() {
	b.handlers = []HandlerTuple{}
	b.imHandlers = []HandlerTuple{}

	for _, module := range b.modules {
		for _, command := range module.GetCommands() {
			for _, rgx := range command.GetMatchers() {
				b.handlers = append(b.handlers, HandlerTuple{
					Pattern: regexp.MustCompile(rgx),
					Handler: command.GetHandler(),
				})
			}

			for _, rgx := range command.GetIMMatchers() {
				b.imHandlers = append(b.imHandlers, HandlerTuple{
					Pattern: regexp.MustCompile(rgx),
					Handler: command.GetHandler(),
				})
			}

			for _, rgx := range command.GetMentionMatchers() {
				b.handlers = append(b.handlers, HandlerTuple{
					Pattern: b.mentionRegex(rgx),
					Handler: command.GetHandler(),
				})
			}
		}
	}
}

// Excuse comes up with an excuse of why something failed.
func (b *SlackBot) Excuse(m messages.Message, err error) {
	b.Slacker.Logger.Error(err)

	if b.Slacker.Config.GetBool("server.serious") {
		b.Slacker.SimplePost(
			m.GetChannel(),
			"An error ocurred and I was unable to fulfill your request.",
			messages.IconDefault,
			true,
		)

		return
	}

	excuses := []string{
		"There is some interference right now and I can't fulfill your request.",
		"Oh noes. I messed up.",
		"Whoops. Something went wrong.",
		"1000s lines of code and I still cant get some things right.",
		"[explodes]",
		"Error: WHY U NO WORK?!",
		"OMG! It failed.",
		"such failure. such request.",
		"Oops I did it again...",
		"A cat is walking over my keywpdfahsgasgdadfk kj h",
	}

	b.Slacker.SimplePost(
		m.GetChannel(),
		excuses[rand.Intn(len(excuses))],
		messages.IconDefault,
		true,
	)
}

// ProcessIMOpen handles IM open events.
func (b *SlackBot) ProcessIMOpen(ev *slack.IMOpenEvent) {
	b.imChannelIDs[ev.Channel] = true
}

// ProcessMessage processes incoming messages events and calls the appropriate
// handlers.
func (b *SlackBot) ProcessMessage(ev *slack.MessageEvent) {
	// Ignore messages from the bot itself.
	if ev.User == b.slackInfo.User.ID {
		return
	}

	// If the message is an IM, use IM handlers.
	if _, ok := b.imChannelIDs[ev.Channel]; ok {
		handled := false

		for _, tuple := range b.imHandlers {
			if result := tuple.Pattern.FindStringSubmatch(ev.Text); result != nil {
				go tuple.Handler(b, messages.NewSlackMessage(ev), result)

				handled = true
			}
		}

		// On an IM, we will show a small help message if no handlers are found.
		if handled == false {
			go b.HandleUsage(messages.NewSlackMessage(ev), []string{})
		}

		return
	}

	for _, tuple := range b.handlers {
		if result := tuple.Pattern.FindStringSubmatch(ev.Text); result != nil {
			go tuple.Handler(b, messages.NewSlackMessage(ev), result)
		}
	}
}

// PostOnFeed posts a message on the feed.
func (b *SlackBot) PostOnFeed(message string) {
	b.Slacker.FeedPost(message)
}

// Post posts a simple messsage to the a channel.
func (b *SlackBot) Post(
	channel string,
	message string,
	icon messages.Icon,
	asUser bool,
) {
	b.Slacker.SimplePost(channel, message, icon, asUser)
}

// PostImage posts a simple message with an image to the channel.
func (b *SlackBot) PostImage(
	channel string,
	message string,
	imageURL string,
	icon messages.Icon,
	asUser bool,
) {
	b.Slacker.SimpleImagePost(channel, message, imageURL, icon, asUser)
}

// GetModules returns the modules used in this bot.
func (b *SlackBot) GetModules() []interfaces.Module {
	return b.modules
}

// StartTyping notify Slack that the bot is "typing".
func (b *SlackBot) StartTyping(channel string) {
	b.slackRTM.SendMessage(b.slackRTM.NewTypingMessage(channel))
}

// GetGonduit gets an instance of a gonduit client.
func (b *SlackBot) GetGonduit() (*gonduit.Conn, error) {
	return b.Slacker.Factory.Make()
}

// GetConfig returns an instance of the configuration store.
func (b *SlackBot) GetConfig() *confer.Config {
	return b.Slacker.Config
}

func (b *SlackBot) GetSlack() *slack.Client {
	return b.Slacker.Slack
}

// HandleUsage shows usage tip.
func (b *SlackBot) HandleUsage(m messages.Message, matches []string) {
	b.Slacker.SimplePost(
		m.GetChannel(),
		"Hi. For usage information, type `help`.",
		messages.IconTasks,
		true,
	)
}

func (b *SlackBot) GetUsername(userId string) (string, error) {
	userInfo, err := b.GetSlack().GetUserInfo(userId)
	if err != nil {
		return "", err
	}

	return userInfo.Name, nil
}