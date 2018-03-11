package constants

// DifferentialQueryCommitHashType are the commit types.
type DifferentialQueryCommitHashType string

const (
	// DifferentialQueryGitCommit is a git commit.
	DifferentialQueryGitCommit DifferentialQueryCommitHashType = "gtcm"
	// DifferentialQueryGitTr is something.
	DifferentialQueryGitTr DifferentialQueryCommitHashType = "gttr"
	// DifferentialQueryHGCommit is a mercurial commit.
	DifferentialQueryHGCommit DifferentialQueryCommitHashType = "hgcm"
)

// DifferentialStatus is the status of a differential revision.
type DifferentialStatus string

const (
	// DifferentialStatusAny is any status.
	DifferentialStatusAny DifferentialStatus = "status-any"
	// DifferentialStatusOpen is any revision that is open.
	DifferentialStatusOpen DifferentialStatus = "status-open"
	// DifferentialStatusAccepted is any revision that is accepted.
	DifferentialStatusAccepted DifferentialStatus = "status-accepted"
	// DifferentialStatusClosed is any revision that is closed.
	DifferentialStatusClosed DifferentialStatus = "status-closed"
)

// DifferentialQueryOrder is the order in which query results cna be ordered.
type DifferentialQueryOrder string

const (
	// DifferentialQueryOrderModified orders results by date modified.
	DifferentialQueryOrderModified DifferentialQueryOrder = "order-modified"
	// DifferentialQueryOrderCreated orders results by date created.
	DifferentialQueryOrderCreated DifferentialQueryOrder = "order-created"
)
