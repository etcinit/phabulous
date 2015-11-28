package bot

import (
	"github.com/nlopes/slack"
)

// BootRTM handles RTM events in Slack.
func (s *SlackService) BootRTM() {
	s.Logger.Infoln("Starting RTM handler...")

	rtm := s.Slack.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			s.Logger.Debugln("RTM Event received.")
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				// Ignore hello

			case *slack.ConnectedEvent:
				s.Logger.Debugln("Infos:", ev.Info)
				s.Logger.Debugln("Connection counter:", ev.ConnectionCount)

				s.Bot = NewBot(s, rtm, ev.Info)

				s.Logger.Infoln("Bot Slack ID: ", ev.Info.User.ID)

				//s.FeedPost("Hi! Phabulous v2 reporting for duty!")
			case *slack.MessageEvent:
				s.Logger.Debugf("Message: %v\n", ev)
				s.Bot.ProcessMessage(ev)

			case *slack.IMOpenEvent:
				s.Bot.ProcessIMOpen(ev)

			case *slack.PresenceChangeEvent:
				s.Logger.Debugf("Presence Change: %v\n", ev)

			case *slack.LatencyReport:
				s.Logger.Debugf("Current latency: %v\n", ev.Value)

			case *slack.RTMError:
				s.Logger.Errorf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				s.Logger.Errorln("Invalid credentials.")
				break Loop

			default:
				// Ignore other events..
			}
		}
	}

	s.Logger.Warnln("RTM handler has stopped.")
}
