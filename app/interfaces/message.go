package interfaces

// Message defines the interface of a message handled by the bot. This allows
// the bot to handle messages from different platforms. Each connector should
// include an implementation of a message for the protocol they connect.
type Message interface {
	// GetChannel returns the channel or room this message was posted on.
	GetChannel() string
	// GetUserID gets the ID or nickname of the user who created this message.
	GetUserID() string
	// GetContent gets the message content.
	GetContent() string
	// GetProviderName returns the name of the provider this message was
	// delivered by. Examples: slack, irc, etc.
	GetProviderName() string
	// IsIM returns true if the message is a direct message sent to the bot.
	IsIM() bool
	// IsSelf returns true if the message was posted by the bot.
	IsSelf() bool
	// HasUser returns true if the message has a user.
	HasUser() bool
}
