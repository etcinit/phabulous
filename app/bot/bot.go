package bot

import (
	"math/rand"
	"regexp"

	"github.com/etcinit/phabulous/app/factories"
	"github.com/etcinit/phabulous/app/messages"
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
	handlers     map[*regexp.Regexp]func(*slack.MessageEvent, []string)
	imHandlers   map[*regexp.Regexp]func(*slack.MessageEvent, []string)
}

func (b *Bot) mentionRegex(contents string) *regexp.Regexp {
	username := b.slackInfo.User.ID

	return regexp.MustCompile("^<@" + username + ">: " + contents + "$")
}

func (b *Bot) loadHandlers() {
	b.handlers = map[*regexp.Regexp]func(*slack.MessageEvent, []string){}
	b.imHandlers = map[*regexp.Regexp]func(*slack.MessageEvent, []string){}

	b.imHandlers[regexp.MustCompile("^lookup ([T|D][0-9]{1,16})$")] =
		b.HandleLookup
	b.imHandlers[regexp.MustCompile("^([T|D][0-9]{1,16})$")] =
		b.HandleLookup
	b.imHandlers[regexp.MustCompile("^help$")] = b.HandleHelp

	b.imHandlers[regexp.MustCompile("dev:post:feed:test")] =
		b.HandleTestFeedMessage

	b.handlers[b.mentionRegex("summon D([0-9]{1,16})")] =
		b.HandleSummon
	b.handlers[b.mentionRegex("help")] = b.HandleHelp
	b.handlers[b.mentionRegex("([T|D][0-9]{1,16})")] = b.HandleLookup
	b.handlers[b.mentionRegex("lookup ([T|D][0-9]{1,16})")] = b.HandleLookup
	b.handlers[regexp.MustCompile("^([T|D][0-9]{1,16})$")] = b.HandleLookup
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
		"Error 500: Internal Server Explosion.",
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

		for re, handler := range b.imHandlers {
			if result := re.FindStringSubmatch(ev.Text); result != nil {
				go handler(ev, result)

				handled = true
			}
		}

		// On an IM, we will show a small help message if no handlers are found.
		if handled == false {
			go b.HandleUsage(ev, []string{})
		}

		return
	}

	for re, handler := range b.handlers {
		if result := re.FindStringSubmatch(ev.Text); result != nil {
			go handler(ev, result)
		}
	}
}
