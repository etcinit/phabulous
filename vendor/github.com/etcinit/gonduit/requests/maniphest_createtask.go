package requests

// ManiphestCreateTaskRequest represents a request to maniphest.createtask.
type ManiphestCreateTaskRequest struct {
	Title        string                        `json:"title"`
	Description  string                        `json:"description"`
	OwnerPHID    string                        `json:"ownerPHID"`
	ViewPolicy   string                        `json:"viewPolicy"`
	EditPolicy   string                        `json:"editPolicy"`
	CCPHIDs      []string                      `json:"ccPHIDs"`
	Priority     int                           `json:"priority"`
	ProjectPHIDs []string                      `json:"projectPHIDs"`
	Request
}
