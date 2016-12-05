package entities

// PhabricatorUser is the partial user information returned by
// phabulous.fromslack.
type PhabricatorUser struct {
	UserPHID      string `json:"userPHID"`
	AccountDomain string `json:"accountDomain"`
}

// SlackUser is the partial information returned by phabulous.toslack.
type SlackUser struct {
	AccountID string `json:"accountID"`
}
