package interfaces

// A Connector provides access to a chat network and supports setting up a Bot
// for interacting with users in that network.
type Connector interface {
	Poster
	Boot() error
	LoadModules(modules []Module)
}
