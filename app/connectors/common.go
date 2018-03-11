package connectors

import (
	"github.com/etcinit/phabulous/app/interfaces"
	"github.com/etcinit/phabulous/app/utilities"
)

// processMessage processes incoming messages events and calls the appropriate
// handlers.
func processMessage(conn interfaces.Bot, msg interfaces.Message) {
	// Ignore messages from the bot itself.
	if msg.IsSelf() || !msg.HasUser() {
		return
	}

	content := msg.GetContent()

	// If the message is an IM, use IM handlers.
	if msg.IsIM() {
		handled := false

		for _, tuple := range conn.GetIMHandlers() {
			var handledResults []string
			pattern := tuple.GetPattern()
			if results := pattern.FindAllStringSubmatch(content, -1); results != nil {

				for _, result := range results {
					result = utilities.UniqueItemsOf(result)
					if !utilities.Contains(handledResults, result) {
						go tuple.GetHandler()(conn, msg, result)
						handledResults = append(handledResults, result...)
						handled = true
					}
				}
			}
		}

		// On an IM, we will show a small help message if no handlers are found.
		if handled == false {
			go conn.GetUsageHandler()(conn, msg, []string{})
		}

		return
	}

	for _, tuple := range conn.GetHandlers() {
		pattern := tuple.GetPattern()
		var handledResults []string

		if results := pattern.FindAllStringSubmatch(content, -1); results != nil {
			for _, result := range results {
				result = utilities.UniqueItemsOf(result)
				if !utilities.Contains(handledResults, result) {
					go tuple.GetHandler()(conn, msg, result)
					handledResults = append(handledResults, result...)
				}
			}
		}
	}
}
