package bot

import (
	"math/rand"
	"regexp"

	"github.com/davecgh/go-spew/spew"
	"github.com/etcinit/gonduit"
	"github.com/etcinit/phabulous/app/factories"
	"github.com/etcinit/phabulous/app/messages"
	"github.com/etcinit/phabulous/app/modules"
	"github.com/etcinit/phabulous/app/modules/core"
	"github.com/etcinit/phabulous/app/modules/dev"
	"github.com/etcinit/phabulous/app/modules/extension"
	"github.com/jacobstr/confer"
	"github.com/nlopes/slack"
)

// NewBot creates a new instance of a Bot.
func NewBot(
	slacker *SlackService,
	slackRTM *slack.RTM,
	slackInfo *slack.Info,
) *Bot {
	bot := &Bot{
		Slacker:      slacker,
		slackInfo:    slackInfo,
		slackRTM:     slackRTM,
		imChannelIDs: map[string]bool{},
		modules: []modules.Module{
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

// Bot represents the state of the bot. It also contains functions related to
// the interactive portion of Phabulous.
type Bot struct {
	Slacker *SlackService
	Factory *factories.GonduitFactory

	slackInfo    *slack.Info
	slackRTM     *slack.RTM
	imChannelIDs map[string]bool
	handlers     []HandlerTuple
	imHandlers   []HandlerTuple

	modules []modules.Module
}

// HandlerTuple a tuples of a pattern and a handler.
type HandlerTuple struct {
	Pattern *regexp.Regexp
	Handler modules.Handler
}

func (b *Bot) mentionRegex(contents string) *regexp.Regexp {
	username := b.slackInfo.User.ID

	return regexp.MustCompile("^<@" + username + ">:? " + contents + "$")
}

func (b *Bot) loadHandlers() {
	b.handlers = []HandlerTuple{}
	b.imHandlers = []HandlerTuple{}

	//b.handlers[b.mentionRegex("summon D([0-9]{1,16})")] =
	//	b.HandleSummon

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

	spew.Dump(b.handlers)
	spew.Dump(b.imHandlers)
}

// Excuse comes up with an excuse of why something failed.
func (b *Bot) Excuse(ev *slack.MessageEvent, err error) {
	b.Slacker.Logger.Error(err)

	if b.Slacker.Config.GetBool("server.serious") {
		b.Slacker.SimplePost(
			ev.Channel,
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
		ev.Channel,
		excuses[rand.Intn(len(excuses))],
		messages.IconDefault,
		true,
	)
}

// ProcessIMOpen handles IM open events.
func (b *Bot) ProcessIMOpen(ev *slack.IMOpenEvent) {
	b.imChannelIDs[ev.Channel] = true
}

// ProcessMessage processes incoming messages events and calls the appropriate
// handlers.
func (b *Bot) ProcessMessage(ev *slack.MessageEvent) {
	// Ignore messages from the bot itself.
	if ev.User == b.slackInfo.User.ID {
		return
	}

	// If the message is an IM, use IM handlers.
	if _, ok := b.imChannelIDs[ev.Channel]; ok {
		handled := false

		for _, tuple := range b.imHandlers {
			if result := tuple.Pattern.FindStringSubmatch(ev.Text); result != nil {
				go tuple.Handler(b, ev, result)

				handled = true
			}
		}

		// On an IM, we will show a small help message if no handlers are found.
		if handled == false {
			go b.HandleUsage(ev, []string{})
		}

		return
	}

	for _, tuple := range b.handlers {
		if result := tuple.Pattern.FindStringSubmatch(ev.Text); result != nil {
			go tuple.Handler(b, ev, result)
		}
	}
}

// PostOnFeed posts a message on the feed.
func (b *Bot) PostOnFeed(message string) {
	b.Slacker.FeedPost(message)
}

// Post posts a simple messsage to the a channel.
func (b *Bot) Post(
	channel string,
	message string,
	icon messages.Icon,
	asUser bool,
) {
	b.Slacker.SimplePost(channel, message, icon, asUser)
}

// PostImage posts a simple message with an image to the channel.
func (b *Bot) PostImage(
	channel string,
	message string,
	imageURL string,
	icon messages.Icon,
	asUser bool,
) {
	b.Slacker.SimpleImagePost(channel, message, imageURL, icon, asUser)
}

// GetModules returns the modules used in this bot.
func (b *Bot) GetModules() []modules.Module {
	return b.modules
}

// StartTyping notify Slack that the bot is "typing".
func (b *Bot) StartTyping(channel string) {
	b.slackRTM.SendMessage(b.slackRTM.NewTypingMessage(channel))
}

// MakeGonduit gets an instance of a gonduit client.
func (b *Bot) MakeGonduit() (*gonduit.Conn, error) {
	return b.Slacker.Factory.Make()
}

// MakeRTM returns an instance of the Slack RTM client.
func (b *Bot) MakeRTM() *slack.RTM {
	return b.slackRTM
}

// MakeSlack returns an instance of the Slack Web client.
func (b *Bot) MakeSlack() *slack.Client {
	return b.Slacker.Slack
}

// MakeConfig returns an instance of the configuration store.
func (b *Bot) MakeConfig() *confer.Config {
	return b.Slacker.Config
}

// HandleUsage shows usage tip.
func (b *Bot) HandleUsage(ev *slack.MessageEvent, matches []string) {
	b.Slacker.SimplePost(
		ev.Channel,
		"Hi. For usage information, type `help`.",
		messages.IconTasks,
		true,
	)
}
