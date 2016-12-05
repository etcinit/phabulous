package interfaces

// A Connector provides access to a chat network and supports setting up a Bot
// for interacting with users in that network.
type Connector interface {
	Poster

	// Boot tells the connector to begin connecting to its network.
	//
	// This method will be called by Phabulous once the server is ready to
	// begin joining networks and all modules have been loaded.
	//
	Boot() error

	// LoadModules provides the connector with a slice of the modules that
	// should be loaded by the connector. The connector should process the
	// modules by creating an internal map of regular expressions and handlers.
	LoadModules(modules []Module)
}
