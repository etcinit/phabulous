package entities

import "github.com/etcinit/gonduit/util"

// DifferentialRevision represents a revision in Differential.
type DifferentialRevision struct {
	ID             string              `json:"id"`
	PHID           string              `json:"phid"`
	Title          string              `json:"title"`
	URI            string              `json:"uri"`
	DateCreated    util.UnixTimestamp  `json:"dateCreated"`
	DateModified   util.UnixTimestamp  `json:"dateModified"`
	AuthorPHID     string              `json:"authorPHID"`
	Status         string              `json:"status"`
	StatusName     string              `json:"statusName"`
	Branch         string              `json:"branch"`
	Summary        string              `json:"summary"`
	TestPlan       string              `json:"testPlan"`
	LineCount      string              `json:"lineCount"`
	ActiveDiffPHID string              `json:"activeDiffPHID"`
	Diffs          []string            `json:"diffs"`
	Commits        []string            `json:"commits"`
	Reviewers      []string            `json:"reviewers"`
	CCs            []string            `json:"ccs"`
	Hashes         [][]string          `json:"hashes"`
	Auxiliary      map[string][]string `json:"auxiliary"`
	RepositoryPHID string              `json:"repositoryPHID"`
}
