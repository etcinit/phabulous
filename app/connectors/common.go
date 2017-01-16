package connectors

import "github.com/etcinit/phabulous/app/interfaces"

// processMessage processes incoming messages events and calls the appropriate
// handlers.
func processMessage(conn interfaces.Bot, msg interfaces.Message) {
	// Ignore messages from the bot itself.
	if msg.IsSelf() {
		return
	}

	content := msg.GetContent()

	// If the message is an IM, use IM handlers.
	if msg.IsIM() {
		handled := false

		for _, tuple := range conn.GetIMHandlers() {
			pattern := tuple.GetPattern()

			if result := pattern.FindStringSubmatch(content); result != nil {
				go tuple.GetHandler()(conn, msg, result)

				handled = true
			}
		}

		// On an IM, we will show a small help message if no handlers are found.
		if handled == false {
			go conn.HandleUsage(msg, []string{})
		}

		return
	}

	for _, tuple := range conn.GetHandlers() {
		pattern := tuple.GetPattern()

		if result := pattern.FindStringSubmatch(content); result != nil {
			go tuple.GetHandler()(conn, msg, result)
		}
	}
}
