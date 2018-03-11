package responses

import "github.com/etcinit/gonduit/entities"

// ManiphestQueryResponse is the response of calling maniphest.query.
type ManiphestQueryResponse map[string]*entities.ManiphestTask

// Get gets the task with the speicfied numeric ID.
func (res ManiphestQueryResponse) Get(key string) *entities.ManiphestTask {
	if _, ok := res[key]; ok {
		return res[key]
	}

	return nil
}
