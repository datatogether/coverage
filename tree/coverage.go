package tree

import (
	"time"
)

// Coverage is information provided by a repository about a url
type Coverage struct {
	// exact url this coverage information is about
	Url string `json:"url,omitempty"`
	// id of the coverage
	RepositoryId string `json:"repositoryId"`
	// sha256 hash of the response (if present)
	Sha256 string `json:"sha256,omitempty"`
	// time of capture from this url
	Timestamp *time.Time `json:"timestamp,omitempty"`
	// url to where this archive now lives
	ArchiveUrl string `json:"archiveUrl,omitempty"`
	// flag for wather or not this archive is in fact complete
	Archived bool `json:"archived"`
	// basic flag for if this url contains data that is difficult
	// to archive with conventional scraping means
	Uncrawlable bool `json:"uncrawlable,omitempty"`
	// a number from 0-10 (10 being high priority) for archiving
	Priority int `json:"priority"`
}
