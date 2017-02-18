package interfaces

// Message defines the interface of a message handled by the bot. This allows
// the bot to handle messages from different platforms. Each connector should
// include an implementation of a message for the protocol they connect.
type Message interface {
	GetChannel() string
	GetUserId() string
	GetContent() string
	GetProviderName() string
	IsIM() bool
	IsSelf() bool
}
