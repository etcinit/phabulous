package responses

import "github.com/etcinit/gonduit/entities"

// ProjectQueryResponse represents a response from calling project.query.
type ProjectQueryResponse struct {
	Data    map[string]entities.Project `json:"data"`
	SlugMap map[string]string           `json:"sligMap"`
	Cursor  entities.Cursor             `json:"cursor"`
}
