package requests

import (
	"github.com/etcinit/gonduit/constants"
	"encoding/json"
)

// DifferentialQueryRequest represents a request to the
// differential.query call.
type DifferentialQueryRequest struct {
	Authors          []string                         `json:"authors"`
	CCs              []string                         `json:"ccs"`
	Reviewers        []string                         `json:"reviewers"`
	Paths            [][]string                       `json:"paths"`
	CommitHashes     [][]string                       `json:"commitHashes"`
	Status           constants.DifferentialStatus     `json:"status"`
	Order            constants.DifferentialQueryOrder `json:"order"`
	Limit            uint64                           `json:"limit"`
	Offset           uint64                           `json:"offset"`
	IDs              []uint64                         `json:"ids"`
	PHIDs            []string                         `json:"phids"`
	Subscribers      []string                         `json:"subscribers"`
	ResponsibleUsers []string                         `json:"responsibleUsers"`
	Branches         []string                         `json:"branches"`
	Request
}

func (dqr *DifferentialQueryRequest) UnmarshalJSON(data []byte) error {
	type Alias DifferentialQueryRequest
	temp := &struct {
		Reviewers map[string]string `json:"reviewers"`
		*Alias
	}{
		Alias: (*Alias)(dqr),
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	dqr.Reviewers = []string{}
	for reviewer := range temp.Reviewers {
		dqr.Reviewers = append(dqr.Reviewers, reviewer)
	}

	return nil
}