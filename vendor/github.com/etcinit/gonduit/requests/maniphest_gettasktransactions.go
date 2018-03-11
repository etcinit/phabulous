package requests

// ManiphestGetTaskTransactions represents a request to maniphest.gettasktransactions.
type ManiphestGetTaskTransactions struct {
	IDs []string `json:"ids"`
	Request
}
