package gerrit

import (
	"context"
	"fmt"
	"net/http"
)

// RevisionKind describes the change kind.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#revision-info
type RevisionKind string

const (
	Rework                 RevisionKind = "REWORK"
	TrivialRebase          RevisionKind = "TRIVIAL_REBASE"
	MergeFirstParentUpdate RevisionKind = "MERGE_FIRST_PARENT_UPDATE"
	NoCodeChange           RevisionKind = "NO_CODE_CHANGE"
	NoChange               RevisionKind = "NO_CHANGE"
)

// WebLinkInfo entity describes a link to an external site.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#web-link-info
type WebLinkInfo struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	ImageURL string `json:"image_url"`
}

// GitPersonInfo entity contains information about the author/committer of a commit.
type GitPersonInfo struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  Timestamp `json:"date"`
	TZ    int       `json:"tz"`
}

// NotifyInfo entity contains detailed information about who should be
// notified about an update
type NotifyInfo struct {
	Accounts []AccountInfo `json:"accounts"`
}

// AbandonInput entity contains information for abandoning a change.
type AbandonInput struct {
	Message       string       `json:"message,omitempty"`
	Notify        string       `json:"notify,omitempty"`
	NotifyDetails []NotifyInfo `json:"notify_details,omitempty"`
}

// ApprovalInfo entity contains information about an approval from a user for a label on a change.
type ApprovalInfo struct {
	AccountInfo
	Value int    `json:"value,omitempty"`
	Date  string `json:"date,omitempty"`
}

// CommitMessageInput entity contains information for changing the commit message of a change.
type CommitMessageInput struct {
	Message       string       `json:"message,omitempty"`
	Notify        string       `json:"notify,omitempty"`
	NotifyDetails []NotifyInfo `json:"notify_details"`
}

// ReadyForReviewInput entity contains information for transitioning a change from WIP to ready.
type ReadyForReviewInput struct {
	Message string `json:"message,omitempty"`
}

// WorkInProgressInput entity contains additional information for a change set to WorkInProgress.
type WorkInProgressInput struct {
	Message string `json:"message,omitempty"`
}

// ChangeEditInput entity contains information for restoring a path within change edit.
type ChangeEditInput struct {
	RestorePath string `json:"restore_path,omitempty"`
	OldPath     string `json:"old_path,omitempty"`
	NewPath     string `json:"new_path,omitempty"`
}

type RestoreChangeEditInput struct {
	RestorePath string `json:"restore_path"`
}

type RenameChangeEditInput struct {
	OldPath string `json:"old_path"`
	NewPath string `json:"new_path"`
}

type PublishChangeEditInput struct {
	Notify string `json:"notify,omitempty"`
}

// ChangeEditMessageInput entity contains information for changing the commit message within a change edit.
type ChangeEditMessageInput struct {
	Message string `json:"message"`
}

// ChangeMessageInfo entity contains information about a message attached to a change.
type ChangeMessageInfo struct {
	ID             string      `json:"id"`
	Author         AccountInfo `json:"author,omitempty"`
	Date           Timestamp   `json:"date"`
	Message        string      `json:"message"`
	Tag            string      `json:"tag,omitempty"`
	RevisionNumber int         `json:"_revision_number,omitempty"`
}

// DeleteChangeMessageInput entity contains the options for deleting a change message.
type DeleteChangeMessageInput struct {
	Reason string `json:"reason,omitempty"`
}

type SubmitRequirementInput struct {
	Name                         string `json:"name"`
	Description                  string `json:"description,omitempty"`
	ApplicabilityExpression      string `json:"applicability_expression,omitempty"`
	SubmittabilityExpression     string `json:"submittability_expression"`
	OverrideExpression           string `json:"override_expression,omitempty"`
	AllowOverrideInChildProjects bool   `json:"allow_override_in_child_projects,omitempt"`
}

type SubmitRequirementResultInfo struct {
	Name                           string                          `json:"name"`
	Description                    string                          `json:"description,omitempty"`
	Status                         string                          `json:"status"`
	IsLegacy                       bool                            `json:"is_legacy"`
	ApplicabilityExpressionResult  SubmitRequirementExpressionInfo `json:"applicability_expression_result,omitempty"`
	SubmittabilityExpressionResult SubmitRequirementExpressionInfo `json:"submittability_expression_result"`
	OverrideExpressionResult       SubmitRequirementExpressionInfo `json:"override_expression_result,omitempty"`
}

type SubmitRequirementExpressionInfo struct {
	Expression   string   `json:"expression,omitempty"`
	Fulfilled    bool     `json:"fulfilled"`
	Status       string   `json:"status"`
	PassingAtoms []string `json:"passing_atoms,omitempty"`
	FailingAtoms []string `json:"failing_atoms,omitempty"`
	ErrorAtoms   []string `json:"error_atoms,omitempty"`
}

// CherryPickInput entity contains information for cherry-picking a change to a new branch.
type CherryPickInput struct {
	Message           string                       `json:"message,omitempty"`
	Destination       string                       `json:"destination"`
	Base              string                       `json:"base,omitempty"`
	Parent            int                          `json:"parent,omitempty"`
	Notify            string                       `json:"notify,omitempty"`
	NotifyDetails     map[RecipientType]NotifyInfo `json:"notify_details,omitempty"`
	KeepReviewers     bool                         `json:"keep_reviewers,omitempty"`
	AllowConflicts    bool                         `json:"allow_conflicts,omitempty"`
	Topic             string                       `json:"topic,omitempty"`
	AllowEmpty        bool                         `json:"allow_empty,omitempty"`
	CommitterEmail    string                       `json:"committer_email,omitempty"`
	ValidationOptions map[string]string            `json:"validation_options,omitempty"`
}

// CommentRange entity describes the range of an inline comment.
type CommentRange struct {
	StartLine      int `json:"start_line"`
	StartCharacter int `json:"start_character"`
	EndLine        int `json:"end_line"`
	EndCharacter   int `json:"end_character"`
}

// DiffFileMetaInfo entity contains meta information about a file diff
type DiffFileMetaInfo struct {
	Name        string        `json:"name"`
	ContentType string        `json:"content_type"`
	Lines       int           `json:"lines"`
	WebLinks    []WebLinkInfo `json:"web_links,omitempty"`
}

// DiffWebLinkInfo entity describes a link on a diff screen to an external site.
type DiffWebLinkInfo struct {
	Name                     string `json:"name"`
	URL                      string `json:"url"`
	ImageURL                 string `json:"image_url"`
	ShowOnSideBySideDiffView bool   `json:"show_on_side_by_side_diff_view"`
	ShowOnUnifiedDiffView    bool   `json:"show_on_unified_diff_view"`
}

// FetchInfo entity contains information about how to fetch a patch set via a certain protocol.
type FetchInfo struct {
	URL      string            `json:"url"`
	Ref      string            `json:"ref"`
	Commands map[string]string `json:"commands,omitempty"`
}

// FixInput entity contains options for fixing commits using the fix change endpoint.
type FixInput struct {
	DeletePatchSetIfCommitMissing bool   `json:"delete_patch_set_if_commit_missing"`
	ExpectMergedAs                string `json:"expect_merged_as"`
}

// PrivateInput entity contains information for changing the private flag on a change.
type PrivateInput struct {
	Message string `json:"message"`
}

// GroupBaseInfo entity contains base information about the group.
type GroupBaseInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// IncludedInInfo entity contains information about the branches a change was merged into and tags it was tagged with.
type IncludedInInfo struct {
	Branches []string          `json:"branches"`
	Tags     []string          `json:"tags"`
	External map[string]string `json:"external,omitempty"`
}

// ProblemInfo entity contains a description of a potential consistency problem with a change.
// These are not related to the code review process, but rather indicate some inconsistency in Gerrit’s database or repository metadata related to the enclosing change.
type ProblemInfo struct {
	Message string `json:"message"`
	Status  string `json:"status,omitempty"`
	Outcome string `json:"outcome,omitempty"`
}

// RebaseInput entity contains information for changing parent when rebasing.
type RebaseInput struct {
	Base               string            `json:"base,omitempty"`
	Strategy           string            `json:"strategy,omitempty"`
	AllowConflicts     bool              `json:"allow_conflicts,omitempty"`
	OnBehalfOfUploader bool              `json:"on_behalf_of_uploader,omitempty"`
	CommitterEmail     string            `json:"committer_email,omitempty"`
	ValidationOptions  map[string]string `json:"validation_options,omitempty"`
}

// RestoreInput entity contains information for restoring a change.
type RestoreInput struct {
	Message string `json:"message,omitempty"`
}

// RevertInput entity contains information for reverting a change.
type RevertInput struct {
	Message string `json:"message,omitempty"`
}

// ReviewInfo entity contains information about a review.
type ReviewInfo struct {
	Labels map[string]int `json:"labels"`
}

// ReviewerUpdateInfo entity contains information about updates
// to change's reviewers set.
type ReviewerUpdateInfo struct {
	Updated   Timestamp   `json:"updated"`    // Timestamp of the update.
	UpdatedBy AccountInfo `json:"updated_by"` // The account which modified state of the reviewer in question.
	Reviewer  AccountInfo `json:"reviewer"`   // The reviewer account added or removed from the change.
	State     string      `json:"state"`      // The reviewer state, one of "REVIEWER", "CC" or "REMOVED".
}

// ReviewResult entity contains information regarding the updates that were
// made to a review.
type ReviewResult struct {
	ReviewInfo
	Reviewers  map[string]ReviewerResult `json:"reviewers,omitempty"`
	Ready      bool                      `json:"ready,omitempty"`
	Error      string                    `json:"error,omitempty"`
	ChangeInfo ChangeInfo                `json:"change_info"`
}

// TopicInput entity contains information for setting a topic.
type TopicInput struct {
	Topic string `json:"topic,omitempty"`
}

// SubmitRecord entity describes results from a submit_rule.
type SubmitRecord struct {
	Status       string                            `json:"status"`
	Ok           map[string]map[string]AccountInfo `json:"ok,omitempty"`
	Reject       map[string]map[string]AccountInfo `json:"reject,omitempty"`
	Need         map[string]interface{}            `json:"need,omitempty"`
	May          map[string]map[string]AccountInfo `json:"may,omitempty"`
	Impossible   map[string]interface{}            `json:"impossible,omitempty"`
	ErrorMessage string                            `json:"error_message,omitempty"`
}

// SubmitInput entity contains information for submitting a change.
type SubmitInput struct {
	OnBehalfOf    string                       `json:"on_behalf_of,omitempty"`
	Notify        string                       `json:"notify,omitempty"`
	NotifyDetails map[RecipientType]NotifyInfo `json:"notify_details,omitempty"`
	WaitForMerge  bool                         `json:"wait_for_merge,omitempty"`
}

// SubmitInfo entity contains information about the change status after submitting.
type SubmitInfo struct {
	Status     string `json:"status"`
	OnBehalfOf string `json:"on_behalf_of,omitempty"`
}

// RuleInput entity contains information to test a Prolog rule.
type RuleInput struct {
	Rule    string `json:"rule"`
	Filters string `json:"filters,omitempty"`
}

// ReviewerInput entity contains information for adding a reviewer to a change.
type ReviewerInput struct {
	Reviewer  string `json:"reviewer"`
	Confirmed bool   `json:"confirmed,omitempty"`
}

// ReviewInput entity contains information for adding a review to a revision.
type ReviewInput struct {
	Message                          string                         `json:"message,omitempty"`
	Tag                              string                         `json:"tag,omitempty"`
	Labels                           map[string]int                 `json:"labels,omitempty"`
	Comments                         map[string][]CommentInput      `json:"comments,omitempty"`
	RobotComments                    map[string][]RobotCommentInput `json:"robot_comments,omitempty"`
	StrictLabels                     bool                           `json:"strict_labels,omitempty"`
	Drafts                           string                         `json:"drafts,omitempty"`
	Notify                           string                         `json:"notify,omitempty"`
	OmitDuplicateComments            bool                           `json:"omit_duplicate_comments,omitempty"`
	OnBehalfOf                       string                         `json:"on_behalf_of,omitempty"`
	Reviewers                        []ReviewerInput                `json:"reviewers,omitempty"`
	Ready                            bool                           `json:"ready,omitempty"`
	WorkInProgress                   bool                           `json:"work_in_progress,omitempty"`
	AddToAttentionSet                []AttentionSetInput            `json:"add_to_attention_set,omitempty"`
	RemoveFromAttentionSet           []AttentionSetInput            `json:"remove_from_attention_set,omitempty"`
	IgnoreAutomaticAttentionSetRules bool                           `json:"ignore_automatic_attention_set_rules,omitempty"`
}

// RelatedChangeAndCommitInfo entity contains information about a related change and commit.
type RelatedChangeAndCommitInfo struct {
	ChangeID              string     `json:"change_id,omitempty"`
	Commit                CommitInfo `json:"commit"`
	ChangeNumber          int        `json:"_change_number,omitempty"`
	RevisionNumber        int        `json:"_revision_number,omitempty"`
	CurrentRevisionNumber int        `json:"_current_revision_number,omitempty"`
	Status                string     `json:"status,omitempty"`
}

// DiffContent entity contains information about the content differences in a file.
type DiffContent struct {
	A      []string          `json:"a,omitempty"`
	B      []string          `json:"b,omitempty"`
	AB     []string          `json:"ab,omitempty"`
	EditA  DiffIntralineInfo `json:"edit_a,omitempty"`
	EditB  DiffIntralineInfo `json:"edit_b,omitempty"`
	Skip   int               `json:"skip,omitempty"`
	Common bool              `json:"common,omitempty"`
}

// CommentInput entity contains information for creating an inline comment.
type CommentInput struct {
	ID         string        `json:"id,omitempty"`
	Path       string        `json:"path,omitempty"`
	Side       string        `json:"side,omitempty"`
	Line       int           `json:"line,omitempty"`
	Range      *CommentRange `json:"range,omitempty"`
	InReplyTo  string        `json:"in_reply_to,omitempty"`
	Updated    *Timestamp    `json:"updated,omitempty"`
	Message    string        `json:"message,omitempty"`
	Unresolved bool          `json:"unresolved,omitempty"`
}

// MoveInput entity contains information for moving a change.
type MoveInput struct {
	DestinationBranch string `json:"destination_branch"`
	Message           string `json:"message,omitempty"`
	KeepAllVotes      bool   `json:"keep_all_votes"`
}

// RobotCommentInput entity contains information for creating an inline robot comment.
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#robot-comment-input
type RobotCommentInput struct {
	CommentInput

	// The ID of the robot that generated this comment.
	RobotID string `json:"robot_id"`
	// An ID of the run of the robot.
	RobotRunID string `json:"robot_run_id"`
	// URL to more information.
	URL string `json:"url,omitempty"`
	// Robot specific properties as map that maps arbitrary keys to values.
	Properties *map[string]*string `json:"properties,omitempty"`
	// Suggested fixes for this robot comment as a list of FixSuggestionInfo
	// entities.
	FixSuggestions []FixSuggestionInfo `json:"fix_suggestions,omitempty"`
}

// RobotCommentInfo entity contains information about a robot inline comment
// RobotCommentInfo has the same fields as CommentInfo. In addition RobotCommentInfo has the following fields:
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#robot-comment-info
type RobotCommentInfo struct {
	CommentInfo

	// The ID of the robot that generated this comment.
	RobotID string `json:"robot_id"`
	// An ID of the run of the robot.
	RobotRunID string `json:"robot_run_id"`
	// URL to more information.
	URL string `json:"url,omitempty"`
	// Robot specific properties as map that maps arbitrary keys to values.
	Properties map[string]string `json:"properties,omitempty"`
	// Suggested fixes for this robot comment as a list of FixSuggestionInfo
	// entities.
	FixSuggestions []FixSuggestionInfo `json:"fix_suggestions,omitempty"`
}

// FixSuggestionInfo entity represents a suggested fix.
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#fix-suggestion-info
type FixSuggestionInfo struct {
	// The UUID of the suggested fix. It will be generated automatically and
	// hence will be ignored if it’s set for input objects.
	FixID string `json:"fix_id"`
	// A description of the suggested fix.
	Description string `json:"description"`
	// A list of FixReplacementInfo entities indicating how the content of one or
	// several files should be modified. Within a file, they should refer to
	// non-overlapping regions.
	Replacements []FixReplacementInfo `json:"replacements"`
}

// FixReplacementInfo entity describes how the content of a file should be replaced by another content.
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#fix-replacement-info
type FixReplacementInfo struct {
	// The path of the file which should be modified. Any file in the repository may be modified.
	Path string `json:"path"`

	// A CommentRange indicating which content of the file should be replaced.
	// Lines in the file are assumed to be separated by the line feed character,
	// the carriage return character, the carriage return followed by the line
	// feed character, or one of the other Unicode linebreak sequences supported
	// by Java.
	Range CommentRange `json:"range"`

	// The content which should be used instead of the current one.
	Replacement string `json:"replacement"`
}

// DiffIntralineInfo entity contains information about intraline edits in a file.
//
// The information consists of a list of <skip length, mark length> pairs,
// where the skip length is the number of characters between the end of
// the previous edit and the start of this edit, and the mark length is the
// number of edited characters following the skip. The start of the edits
// is from the beginning of the related diff content lines.
//
// Note that the implied newline character at the end of each line
// is included in the length calculation, and thus it is possible for
// the edits to span newlines.
type DiffIntralineInfo [][2]int

// ChangeInput entity contains information about creating a new change.
//
// Docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#change-input
type ChangeInput struct {
	Project           string                 `json:"project"`
	Branch            string                 `json:"branch"`
	Subject           string                 `json:"subject"`
	Topic             string                 `json:"topic,omitempty"`
	Status            string                 `json:"status,omitempty"`
	IsPrivate         bool                   `json:"is_private,omitempty"`
	WorkInProgress    bool                   `json:"work_in_progress,omitempty"`
	BaseChange        string                 `json:"base_change,omitempty"`
	BaseCommit        string                 `json:"base_commit,omitempty"`
	NewBranch         bool                   `json:"new_branch,omitempty"`
	ValidationOptions map[string]interface{} `json:"validation_options,omitempty"`
	Merge             *MergeInput            `json:"merge,omitempty"`
	Author            *AccountInput          `json:"author,omitempty"`
	Notify            string                 `json:"notify,omitempty"`
	NotifyDetails     string                 `json:"notify_details,omitempty"`
}

// ChangeInfo entity contains information about a change.
type ChangeInfo struct {
	ID                     string                      `json:"id"`
	URL                    string                      `json:"url,omitempty"`
	Project                string                      `json:"project"`
	Branch                 string                      `json:"branch"`
	Topic                  string                      `json:"topic,omitempty"`
	AttentionSet           map[string]AttentionSetInfo `json:"attention_set,omitempty"`
	Assignee               AccountInfo                 `json:"assignee,omitempty"`
	Hashtags               []string                    `json:"hashtags,omitempty"`
	ChangeID               string                      `json:"change_id"`
	Subject                string                      `json:"subject"`
	Status                 string                      `json:"status"`
	Created                Timestamp                   `json:"created"`
	Updated                Timestamp                   `json:"updated"`
	Submitted              *Timestamp                  `json:"submitted,omitempty"`
	Submitter              AccountInfo                 `json:"submitter,omitempty"`
	Starred                bool                        `json:"starred,omitempty"`
	Reviewed               bool                        `json:"reviewed,omitempty"`
	SubmitType             string                      `json:"submit_type,omitempty"`
	Mergeable              bool                        `json:"mergeable,omitempty"`
	Submittable            bool                        `json:"submittable,omitempty"`
	Insertions             int                         `json:"insertions"`
	Deletions              int                         `json:"deletions"`
	TotalCommentCount      int                         `json:"total_comment_count,omitempty"`
	UnresolvedCommentCount int                         `json:"unresolved_comment_count,omitempty"`
	Number                 int                         `json:"_number"`
	Owner                  AccountInfo                 `json:"owner"`
	Actions                map[string]ActionInfo       `json:"actions,omitempty"`
	Labels                 map[string]LabelInfo        `json:"labels,omitempty"`
	PermittedLabels        map[string][]string         `json:"permitted_labels,omitempty"`
	RemovableReviewers     []AccountInfo               `json:"removable_reviewers,omitempty"`
	Reviewers              map[string][]AccountInfo    `json:"reviewers,omitempty"`
	PendingReviewers       map[string][]AccountInfo    `json:"pending_reviewers,omitempty"`
	ReviewerUpdates        []ReviewerUpdateInfo        `json:"reviewer_updates,omitempty"`
	Messages               []ChangeMessageInfo         `json:"messages,omitempty"`
	CurrentRevision        string                      `json:"current_revision,omitempty"`
	Revisions              map[string]RevisionInfo     `json:"revisions,omitempty"`
	MoreChanges            bool                        `json:"_more_changes,omitempty"`
	Problems               []ProblemInfo               `json:"problems,omitempty"`
	IsPrivate              bool                        `json:"is_private,omitempty"`
	WorkInProgress         bool                        `json:"work_in_progress,omitempty"`
	HasReviewStarted       bool                        `json:"has_review_started,omitempty"`
	RevertOf               int                         `json:"revert_of,omitempty"`
	SubmissionID           string                      `json:"submission_id,omitempty"`
	CherryPickOfChange     int                         `json:"cherry_pick_of_change,omitempty"`
	CherryPickOfPatchSet   int                         `json:"cherry_pick_of_patch_set,omitempty"`
	ContainsGitConflicts   bool                        `json:"contains_git_conflicts,omitempty"`
	BaseChange             string                      `json:"base_change,omitempty"`
}

// LabelInfo entity contains information about a label on a change, always corresponding to the current patch set.
type LabelInfo struct {
	Optional bool `json:"optional,omitempty"`

	// Fields set by LABELS
	Approved     AccountInfo `json:"approved,omitempty"`
	Rejected     AccountInfo `json:"rejected,omitempty"`
	Recommended  AccountInfo `json:"recommended,omitempty"`
	Disliked     AccountInfo `json:"disliked,omitempty"`
	Blocking     bool        `json:"blocking,omitempty"`
	Value        int         `json:"value,omitempty"`
	DefaultValue int         `json:"default_value,omitempty"`

	// Fields set by DETAILED_LABELS
	All    []ApprovalInfo    `json:"all,omitempty"`
	Values map[string]string `json:"values,omitempty"`
}

// The MergeInput entity contains information about the merge
//
// Docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#merge-input
type MergeInput struct {
	Source         string `json:"source"`
	SourceBranch   string `json:"source_branch,omitempty"`
	Strategy       string `json:"strategy,omitempty"`
	AllowConflicts bool   `json:"allow_conflicts,omitempty"`
}

// The ParentInfo entity contains information about the parent commit of a patch-set.
type ParentInfo struct {
	BranchName             string `json:"branch_name,omitempty"`
	CommitID               string `json:"commit_id,omitempty"`
	IsMergedInTargetBranch bool   `json:"is_merged_in_target_branch"`
	ChangeID               string `json:"change_id,omitempty"`
	ChangeNumber           int    `json:"change_number,omitempty"`
	PatchSetNumber         int    `json:"patch_set_number,omitempty"`
	ChangeStatus           string `json:"change_status,omitempty"`
}

// RevisionInfo entity contains information about a patch set.
type RevisionInfo struct {
	Kind              RevisionKind          `json:"kind,omitempty"`
	Draft             bool                  `json:"draft,omitempty"`
	Number            int                   `json:"_number"`
	Created           Timestamp             `json:"created"`
	Uploader          AccountInfo           `json:"uploader"`
	Ref               string                `json:"ref"`
	Fetch             map[string]FetchInfo  `json:"fetch"`
	Commit            CommitInfo            `json:"commit,omitempty"`
	Files             map[string]FileInfo   `json:"files,omitempty"`
	Actions           map[string]ActionInfo `json:"actions,omitempty"`
	Reviewed          bool                  `json:"reviewed,omitempty"`
	MessageWithFooter string                `json:"messageWithFooter,omitempty"`
	ParentsData       []ParentInfo          `json:"parents_data,omitempty"`
}

// CommentInfo entity contains information about an inline comment.
type CommentInfo struct {
	PatchSet        int           `json:"patch_set,omitempty"`
	ID              string        `json:"id"`
	Path            string        `json:"path,omitempty"`
	Side            string        `json:"side,omitempty"`
	Line            int           `json:"line,omitempty"`
	Range           *CommentRange `json:"range,omitempty"`
	InReplyTo       string        `json:"in_reply_to,omitempty"`
	Message         string        `json:"message,omitempty"`
	Updated         *Timestamp    `json:"updated"`
	Author          AccountInfo   `json:"author,omitempty"`
	Unresolved      *bool         `json:"unresolved,omitempty"`
	ChangeMessageID string        `json:"change_message_id,omitempty"`
	CommitID        string        `json:"commit_id,omitempty"`
}

// QueryOptions specifies global parameters to query changes / reviewers.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-changes
type QueryOptions struct {
	// Query parameter
	// Clients are allowed to specify more than one query by setting the q parameter multiple times.
	// In this case the result is an array of arrays, one per query in the same order the queries were given in.
	//
	// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/user-search.html#_search_operators
	Query []string `url:"q,omitempty"`

	// The n parameter can be used to limit the returned results.
	// If the n query parameter is supplied and additional changes exist that match the query beyond the end, the last change object has a _more_changes: true JSON field set.
	Limit int `url:"n,omitempty"`

	// The S or start query parameter can be supplied to skip a number of changes from the list.
	Start int `url:"start,omitempty"`
}

// QueryChangeOptions specifies the parameters to the ChangesService.Query.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-changes
type QueryChangeOptions struct {
	QueryOptions

	// The S or start query parameter can be supplied to skip a number of changes from the list.
	Skip  int `url:"S,omitempty"`
	Start int `url:"start,omitempty"`

	ChangeOptions
}

// ChangeOptions specifies the parameters for Query changes.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-changes
type ChangeOptions struct {
	// Additional fields can be obtained by adding o parameters, each option requires more database lookups and slows down the query response time to the client so they are generally disabled by default.
	//
	// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-changes
	AdditionalFields []string `url:"o,omitempty"`
}

type Change struct {
	Raw    *ChangeInfo
	gerrit *Gerrit
	Base   string
}

type ChangeService struct {
	gerrit *Gerrit
}

// Query lists changes visible to the caller.
// The query string must be provided by the q parameter.
// The n parameter can be used to limit the returned results.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-changes
func (s *ChangeService) Query(ctx context.Context, opt *QueryChangeOptions) (*[]ChangeInfo, *http.Response, error) {
	v := new([]ChangeInfo)
	resp, err := s.gerrit.Requester.Call(ctx, "GET", "changes/", opt, v)
	return v, resp, err
}

// Get retrieves a change.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-change
func (s *ChangeService) Get(ctx context.Context, changeID string, AdditionalFields ...string) (*Change, *http.Response, error) {
	change := Change{Raw: new(ChangeInfo), gerrit: s.gerrit, Base: changeID}

	opt := new(ChangeOptions)
	opt.AdditionalFields = append(opt.AdditionalFields, AdditionalFields...)

	resp, err := change.Poll(ctx, opt)
	if err != nil {
		return nil, resp, err
	}

	return &change, resp, nil
}

// Create creates a new change.
// The change input ChangeInput entity must be provided in the request body.
// As response a ChangeInfo entity is returned that describes the resulting change.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#create-change
func (s *ChangeService) Create(ctx context.Context, input *ChangeInput) (*Change, *http.Response, error) {
	obj := Change{Raw: new(ChangeInfo), gerrit: s.gerrit, Base: ""}
	return obj.Create(ctx, input)
}

// Delete deletes a new or abandoned change
// New or abandoned changes can be deleted by their owner if the user is granted the Delete Own Changes
// permission, otherwise only by administrators.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#delete-change
func (s *ChangeService) Delete(ctx context.Context, changeID string) (bool, *http.Response, error) {
	obj := Change{Raw: new(ChangeInfo), gerrit: s.gerrit, Base: changeID}
	return obj.Delete(ctx)
}

func (c *Change) Poll(ctx context.Context, opt *ChangeOptions) (*http.Response, error) {
	u := fmt.Sprintf("changes/%s", c.Base)
	return c.gerrit.Requester.Call(ctx, "GET", u, opt, c.Raw)
}

func (c *Change) Create(ctx context.Context, input *ChangeInput) (*Change, *http.Response, error) {
	v := new(ChangeInfo)
	resp, err := c.gerrit.Requester.Call(ctx, "POST", "changes/", input, v)

	if err != nil {
		return nil, resp, err
	}

	c.Base = v.ID
	resp, err = c.Poll(ctx, nil)
	if err != nil {
		return nil, resp, err
	}
	return c, resp, nil
}

func (c *Change) Delete(ctx context.Context) (bool, *http.Response, error) {
	u := fmt.Sprintf("changes/%s", c.Base)
	resp, err := c.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)

	if err != nil {
		return false, resp, err
	}
	return true, resp, nil
}

// GetDetail retrieves a change with labels, detailed labels, detailed accounts, and messages.
//
// This response will contain all votes for each label and include one combined vote.
// The combined label vote is calculated in the following order (from highest to lowest): REJECTED > APPROVED > DISLIKED > RECOMMENDED.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-change-detail
func (c *Change) GetDetail(ctx context.Context, opt *ChangeOptions) (*ChangeInfo, *http.Response, error) {
	u := fmt.Sprintf("changes/%s/detail", c.Base)

	v := new(ChangeInfo)
	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, opt, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// SetCommitMessage creates a new patch set with a new commit message.
// The new commit message must be provided in the request body inside a CommitMessageInput entity.
// If a Change-Id footer is specified, it must match the current Change-Id footer.
// If the Change-Id footer is absent, the current Change-Id is added to the message.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#set-message
func (c *Change) SetCommitMessage(ctx context.Context, input *CommitMessageInput) (bool, *http.Response, error) {
	u := fmt.Sprintf("changes/%s/message", c.Base)

	resp, err := c.gerrit.Requester.Call(ctx, "PUT", u, input, nil)
	if err != nil {
		return false, resp, err
	}
	return true, resp, nil
}

// SetReadyForReview marks the change as ready for review (set WIP property to false)
// Changes may only be marked ready by the owner, project owners or site administrators.
// Activates notifications of reviewer. The request body does not need to include a
// WorkInProgressInput entity if no review comment is added.
// Marking a change ready for review also adds all of the reviewers of the change to the attention set.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#set-ready-for-review
func (c *Change) SetReadyForReview(ctx context.Context, input *ReadyForReviewInput) (bool, *http.Response, error) {
	u := fmt.Sprintf("changes/%s/ready", c.Base)

	resp, err := c.gerrit.Requester.Call(ctx, "POST", u, input, nil)
	if err != nil {
		return false, resp, err
	}
	return true, resp, nil
}

// SetWorkInProgress Marks the change as not ready for review yet.
// Changes may only be marked not ready by the owner, project owners or site administrators.
// The request body does not need to include a WorkInProgressInput entity if no review comment is added.
// Actions that create a new patch set in a WIP change default to notifying OWNER instead of ALL.
// Marking a change work in progress also removes all users from the attention set.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#set-work-in-pogress
func (c *Change) SetWorkInProgress(ctx context.Context, input *WorkInProgressInput) (bool, *http.Response, error) {
	u := fmt.Sprintf("changes/%s/wip", c.Base)

	resp, err := c.gerrit.Requester.Call(ctx, "POST", u, input, nil)
	if err != nil {
		return false, resp, err
	}
	return true, resp, nil
}

// GetTopic retrieves the topic of a change.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-topic
func (c *Change) GetTopic(ctx context.Context) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("changes/%s/topic", c.Base)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// SetTopic sets the topic of a change.
// The new topic must be provided in the request body inside a TopicInput entity.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#set-topic
func (c *Change) SetTopic(ctx context.Context, input *TopicInput) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("changes/%s/topic", c.Base)

	resp, err := c.gerrit.Requester.Call(ctx, "PUT", u, input, v)

	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// DeleteTopic deletes the topic of a change.
// The request body does not need to include a TopicInput entity if no review comment is added.
//
// Please note that some proxies prohibit request bodies for DELETE requests.
// In this case, if you want to specify a commit message, use PUT to delete the topic.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#delete-topic
func (c *Change) DeleteTopic(ctx context.Context) (bool, *http.Response, error) {
	u := fmt.Sprintf("changes/%s/topic", c.Base)
	resp, err := c.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)

	if err != nil {
		return false, resp, err
	}
	return true, resp, nil
}

// Abandon abandons a change.
//
// The request body does not need to include a AbandonInput entity if no review
// comment is added.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#abandon-change
func (c *Change) Abandon(ctx context.Context, input *AbandonInput) (*ChangeInfo, *http.Response, error) {
	return c.operate(ctx, "abandon", input)
}

// Restore Restores a change.
//
// The request body does not need to include a RestoreInput entity if no review
// comment is added.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#restore-change
func (c *Change) Restore(ctx context.Context, input *RestoreInput) (*ChangeInfo, *http.Response, error) {
	return c.operate(ctx, "restore", input)
}

// Rebase Rebases a change.
//
// Optionally, the parent revision can be changed to another patch set through the RebaseInput entity.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#rebase-change
func (c *Change) Rebase(ctx context.Context, input *RebaseInput) (*ChangeInfo, *http.Response, error) {
	return c.operate(ctx, "rebase", input)
}

// Move moves a change.
//
// The destination branch must be provided in the request body inside a MoveInput entity.
// Only veto votes that are blocking the change from submission are moved to the destination branch by default.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#move-change
func (c *Change) Move(ctx context.Context, input *MoveInput) (*ChangeInfo, *http.Response, error) {
	return c.operate(ctx, "move", input)
}

// Revert Reverts a change.
//
// The request body does not need to include a RevertInput entity if no review
// comment is added.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#revert-change
func (c *Change) Revert(ctx context.Context, input *RevertInput) (*ChangeInfo, *http.Response, error) {
	return c.operate(ctx, "revert", input)
}

// Submit Submits a change.
//
// The request body only needs to include a SubmitInput entity if submitting on behalf of another user.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#submit-change
func (c *Change) Submit(ctx context.Context, input *SubmitInput) (*ChangeInfo, *http.Response, error) {
	return c.operate(ctx, "submit", input)
}

// Fix performs consistency checks on the change as with GET /check, and additionally fixes any problems that can be fixed automatically.
// The returned field values reflect any fixes.
//
// Some fixes have options controlling their behavior, which can be set in the FixInput entity body.
// Only the change owner, a project owner, or an administrator may fix changes.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#fix-change
func (c *Change) Fix(ctx context.Context, input *FixInput) (*ChangeInfo, *http.Response, error) {
	return c.operate(ctx, "check", input)
}

// MarkPrivate Marks the change to be private. Only open changes can be marked private.
// Changes may only be marked private by the owner or site administrators.
// A message can be specified in the request body inside a PrivateInput entity.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#mark-private
func (c *Change) MarkPrivate(ctx context.Context, input *PrivateInput) (bool, *http.Response, error) {
	u := fmt.Sprintf("changes/%s/private", c.Base)

	resp, err := c.gerrit.Requester.Call(ctx, "POST", u, input, nil)
	if err != nil {
		return false, resp, err
	}
	return true, resp, nil
}

// UnmarkPrivate Marks the change to be non-private. Note users can only unmark own private changes.
// If the change was already not private, the response is “409 Conflict”.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#unmark-private
func (c *Change) UnmarkPrivate(ctx context.Context) (bool, *http.Response, error) {
	u := fmt.Sprintf("changes/%s/private", c.Base)

	resp, err := c.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)
	if err != nil {
		return false, resp, err
	}
	return true, resp, nil
}

// operate is an internal function to consolidate code used by SubmitChange,
// AbandonChange and other similar functions.
func (c *Change) operate(ctx context.Context, tail string, input interface{}) (*ChangeInfo, *http.Response, error) {
	v := new(ChangeInfo)
	u := fmt.Sprintf("changes/%s/%s", c.Base, tail)

	resp, err := c.gerrit.Requester.Call(ctx, "POST", u, input, v)

	if err != nil {
		return nil, resp, err
	}

	return v, resp, nil
}

// SubmittedTogether returns a list of all changes which are submitted when {submit} is called for this change, including the current change itself.
// An empty list is returned if this change will be submitted by itself (no other changes).
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#submitted_together
func (c *Change) SubmittedTogether(ctx context.Context) (*[]ChangeInfo, *http.Response, error) {
	u := fmt.Sprintf("changes/%s/submitted_together", c.Base)

	v := new([]ChangeInfo)
	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetIncludedIn retrieves the branches and tags in which a change is included.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-included-in
func (c *Change) GetIncludedIn(ctx context.Context) (*IncludedInInfo, *http.Response, error) {
	v := new(IncludedInInfo)
	u := fmt.Sprintf("changes/%s/in", c.Base)

	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// ListComments lists the published comments of all revisions of the change.
// The entries in the map are sorted by file path, and the comments for each path are sorted by patch set number.
// Each comment has the patch_set and author fields set.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-change-comments
func (c *Change) ListComments(ctx context.Context) (map[string][]CommentInfo, *http.Response, error) {
	v := make(map[string][]CommentInfo)
	u := fmt.Sprintf("changes/%s/comments", c.Base)
	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, &v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// ListDrafts lLists the draft comments of all revisions of the change that belong to the calling user.
// The entries in the map are sorted by file path, and the comments for each path are sorted by patch set number.
// Each comment has the patch_set field set, and no author.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-change-drafts
func (c *Change) ListDrafts(ctx context.Context) (map[string][]CommentInfo, *http.Response, error) {
	v := make(map[string][]CommentInfo)
	u := fmt.Sprintf("changes/%s/drafts", c.Base)
	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, &v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// Check performs consistency checks on the change, and returns a ChangeInfo entity with the problems field set to a list of ProblemInfo entities.
// Depending on the type of problem, some fields not marked optional may be missing from the result.
// At least id, project, branch, and _number will be present.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#check-change
func (c *Change) Check(ctx context.Context) (*ChangeInfo, *http.Response, error) {
	v := new(ChangeInfo)
	u := fmt.Sprintf("changes/%s/check", c.Base)
	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return nil, resp, err
	}

	return v, resp, nil
}

// Index adds or updates the change in the secondary index.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#index-change
func (c *Change) Index(ctx context.Context) (*http.Response, error) {
	u := fmt.Sprintf("changes/%s/index", c.Base)
	return c.gerrit.Requester.Call(ctx, "POST", u, nil, nil)
}

// GetHashtags gets the hashtags associated with a change.
//
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-hashtags
func (c *Change) GetHashtags(ctx context.Context) ([]string, *http.Response, error) {
	v := new([]string)
	u := fmt.Sprintf("changes/%s/hashtags", c.Base)
	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return nil, resp, err
	}

	return *v, resp, nil
}

// SetHashtags adds and/or removes hashtags from a change.
//
// As response the change’s hashtags are returned as a list of strings.
//
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#set-hashtags
func (c *Change) SetHashtags(ctx context.Context, input *HashtagsInput) ([]string, *http.Response, error) {
	v := new([]string)
	u := fmt.Sprintf("changes/%s/hashtags", c.Base)
	resp, err := c.gerrit.Requester.Call(ctx, "POST", u, input, v)

	if err != nil {
		return nil, resp, err
	}

	return *v, resp, nil
}

// ListMessages lists all the messages of a change including detailed account information.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-change-messages
func (c *Change) ListMessages(ctx context.Context) (*[]ChangeMessageInfo, *http.Response, error) {
	v := new([]ChangeMessageInfo)
	u := fmt.Sprintf("changes/%s/messages", c.Base)
	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetMessage Retrieves a change message including detailed account information.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-change-message
func (c *Change) GetMessage(ctx context.Context, messageID string) (*ChangeMessageInfo, *http.Response, error) {
	v := new(ChangeMessageInfo)
	u := fmt.Sprintf("changes/%s/messages/%s", c.Base, messageID)
	resp, err := c.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// DeleteMessage Deletes a change message by replacing the change message with a new message,
// which contains the name of the user who deleted the change message and the reason why it was deleted.
//
//        The reason can be provided in the request body as a DeleteChangeMessageInput entity.
//
// Note that only users with the Administrate Server global capability are permitted to delete a change message.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#delete-change-message
func (c *Change) DeleteMessage(ctx context.Context, messageID string, input *DeleteChangeMessageInput) (*ChangeMessageInfo, *http.Response, error) {
	v := new(ChangeMessageInfo)
	u := fmt.Sprintf("changes/%s/messages/%s/delete", c.Base, messageID)
	resp, err := c.gerrit.Requester.Call(ctx, "POST", u, input, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// CheckSubmitRequirements a submit requirement and returns the result as a SubmitRequirementResultInfo.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#check-submit-requirements
func (c *Change) CheckSubmitRequirements(ctx context.Context, input *SubmitRequirementInput) (*SubmitRequirementResultInfo, *http.Response, error) {
	v := new(SubmitRequirementResultInfo)
	u := fmt.Sprintf("changes/%s/check.submit_requirement", c.Base)
	resp, err := c.gerrit.Requester.Call(ctx, "POST", u, input, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}
