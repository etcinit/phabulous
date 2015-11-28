package messages

import "github.com/etcinit/gonduit/constants"

// Icon is the type of the PHID.
type Icon string

const (
	// IconDefault is used for regular messages.
	IconDefault Icon = "http://i.imgur.com/7Hzgo9Y.png"

	// IconCommits is used for commit-related messages.
	IconCommits Icon = "http://i.imgur.com/v8ReRKx.png"

	// IconTasks is used for task-related messages.
	IconTasks Icon = "http://i.imgur.com/jD7rf9x.png"

	// IconRevisions is used for revision-related messages.
	IconRevisions Icon = "http://i.imgur.com/NiPouYj.png"
)

// PhidTypeToIcon gets the matching icon for a PHID type.
func PhidTypeToIcon(phidType constants.PhidType) Icon {
	switch phidType {
	case constants.PhidTypeCommit:
		return IconCommits
	case constants.PhidTypeTask:
		return IconTasks
	case constants.PhidTypeDifferentialRevision:
		return IconRevisions
	default:
		return IconDefault
	}
}
