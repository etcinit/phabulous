package utilities

import "github.com/etcinit/phabulous/app/interfaces"

// HasModule returns whether or not a bot connector has a specific module
// loaded.
func HasModule(bot interfaces.Bot, name string) bool {
	for _, module := range bot.GetModules() {
		if module.GetName() == name {
			return true
		}
	}

	return false
}
