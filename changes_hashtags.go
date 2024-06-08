package gerrit

// HashtagsInput entity contains information about hashtags to add to, and/or remove from, a change.
//
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#hashtags-input
type HashtagsInput struct {
	// The list of hashtags to be added to the change.
	Add []string `json:"add,omitempty"`

	// The list of hashtags to be removed from the change.
	Remove []string `json:"remove,omitempty"`
}
