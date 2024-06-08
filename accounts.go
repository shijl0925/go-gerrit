package gerrit

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

// AccountsService contains Account related REST endpoints
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html
type AccountsService struct {
	gerrit *Gerrit
}

type Account struct {
	Raw    *AccountInfo
	gerrit *Gerrit
	Base   string
}

// AccountInfo entity contains information about an account.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#account-info
type AccountInfo struct {
	AccountID   int    `json:"_account_id,omitempty"`
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Email       string `json:"email,omitempty"`
	Username    string `json:"username,omitempty"`

	// Avatars lists avatars of various sizes for the account.
	// This field is only populated if the avatars plugin is enabled.
	Avatars []struct {
		URL    string `json:"url,omitempty"`
		Height int    `json:"height,omitempty"`
	} `json:"avatars,omitempty"`
	MoreAccounts    bool     `json:"_more_accounts,omitempty"`
	SecondaryEmails []string `json:"secondary_emails,omitempty"`
	Status          string   `json:"status,omitempty"`
	Inactive        bool     `json:"inactive,omitempty"`
	Tags            []string `json:"tags,omitempty"`
}

// QueryAccountOptions queries accounts visible to the caller.
type QueryAccountOptions struct {
	QueryOptions

	// The `S` or `start` query parameter can be supplied to skip a number of accounts from the list.
	Start int `url:"S,omitempty"`

	AccountOptions
}

// AccountOptions specifies parameters for Query Accounts.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#query-account
type AccountOptions struct {
	// Additional fields can be obtained by adding o parameters.
	// Currently supported are "DETAILS" and "ALL_EMAILS".
	//
	// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#query-account
	AdditionalFields []string `url:"o,omitempty"`
}

// SSHKeyInfo entity contains information about an SSH key of a user.
type SSHKeyInfo struct {
	Seq          int    `json:"seq"`
	SSHPublicKey string `json:"ssh_public_key"`
	EncodedKey   string `json:"encoded_key"`
	Algorithm    string `json:"algorithm"`
	Comment      string `json:"comment,omitempty"`
	Valid        bool   `json:"valid"`
}

// UsernameInput entity contains information for setting the username for an account.
type UsernameInput struct {
	Username string `json:"username"`
}

// QueryLimitInfo entity contains information about the Query Limit of a user.
type QueryLimitInfo struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// HTTPPasswordInput entity contains information for setting/generating an HTTP password.
type HTTPPasswordInput struct {
	Generate     bool   `json:"generate,omitempty"`
	HTTPPassword string `json:"http_password,omitempty"`
}

type OAuthTokenInfo struct {
	Username     string `json:"username"`
	ResourceHost string `json:"resource_host"`
	AccessToken  string `json:"access_token"`
	ProviderID   string `json:"provider_id,omitempty"`
	ExpiresAt    string `json:"expires_at,omitempty"`
	Type         string `json:"type"`
}

// GpgKeysInput entity contains information for adding/deleting GPG keys.
type GpgKeysInput struct {
	Add    []string `json:"add"`
	Delete []string `json:"delete"`
}

// GpgKeyInfo entity contains information about a GPG public key.
type GpgKeyInfo struct {
	ID          string   `json:"id,omitempty"`
	Fingerprint string   `json:"fingerprint,omitempty"`
	UserIDs     []string `json:"user_ids,omitempty"`
	Key         string   `json:"key,omitempty"`
}

// EmailInput entity contains information for registering a new email address.
type EmailInput struct {
	Email          string `json:"email"`
	Preferred      bool   `json:"preferred,omitempty"`
	NoConfirmation bool   `json:"no_confirmation,omitempty"`
}

// EmailInfo entity contains information about an email address of a user.
type EmailInfo struct {
	Email               string `json:"email"`
	Preferred           bool   `json:"preferred,omitempty"`
	PendingConfirmation bool   `json:"pending_confirmation,omitempty"`
}

// AccountInput entity contains information for the creation of a new account.
type AccountInput struct {
	Username     string   `json:"username,omitempty"`
	Name         string   `json:"name,omitempty"`
	Email        string   `json:"email,omitempty"`
	SSHKey       string   `json:"ssh_key,omitempty"`
	HTTPPassword string   `json:"http_password,omitempty"`
	Groups       []string `json:"groups,omitempty"`
}

// AccountDetailInfo entity contains detailed information about an account.
type AccountDetailInfo struct {
	AccountInfo
	RegisteredOn Timestamp `json:"registered_on"`
}

// AccountExternalIdInfo entity contains information for an external id of an account.
type AccountExternalIdInfo struct {
	Identity     string `json:"identity"`
	EmailAddress string `json:"email_address,omitempty"`
	Trusted      bool   `json:"trusted"`
	CanDelete    bool   `json:"can_delete,omitempty"`
}

// AccountNameInput entity contains information for setting a name for an account.
type AccountNameInput struct {
	Name string `json:"name,omitempty"`
}

// AccountStatusInput entity contains information for setting a status for an account.
type AccountStatusInput struct {
	Status string `json:"status,omitempty"` //The new status of the account. If not set or if set to an empty string, the account status is deleted.
}

type DisplayNameInput struct {
	DisplayName string `json:"display_name"`
}

// AccountCapabilityInfo entity contains information about the global capabilities of a user.
type AccountCapabilityInfo struct {
	AccessDatabase     bool           `json:"accessDatabase,omitempty"`
	AdministrateServer bool           `json:"administrateServer,omitempty"`
	CreateAccount      bool           `json:"createAccount,omitempty"`
	CreateGroup        bool           `json:"createGroup,omitempty"`
	CreateProject      bool           `json:"createProject,omitempty"`
	EmailReviewers     bool           `json:"emailReviewers,omitempty"`
	FlushCaches        bool           `json:"flushCaches,omitempty"`
	KillTask           bool           `json:"killTask,omitempty"`
	MaintainServer     bool           `json:"maintainServer,omitempty"`
	Priority           string         `json:"priority,omitempty"`
	QueryLimit         QueryLimitInfo `json:"queryLimit"`
	RunAs              bool           `json:"runAs,omitempty"`
	RunGC              bool           `json:"runGC,omitempty"`
	StreamEvents       bool           `json:"streamEvents,omitempty"`
	ViewAllAccounts    bool           `json:"viewAllAccounts,omitempty"`
	ViewCaches         bool           `json:"viewCaches,omitempty"`
	ViewConnections    bool           `json:"viewConnections,omitempty"`
	ViewPlugins        bool           `json:"viewPlugins,omitempty"`
	ViewQueue          bool           `json:"viewQueue,omitempty"`
}

// DiffPreferencesInfo entity contains information about the diff preferences of a user.
type DiffPreferencesInfo struct {
	Context                 int    `json:"context"`
	Theme                   string `json:"theme"`
	ExpandAllComments       bool   `json:"expand_all_comments,omitempty"`
	IgnoreWhitespace        string `json:"ignore_whitespace"`
	IntralineDifference     bool   `json:"intraline_difference,omitempty"`
	LineLength              int    `json:"line_length"`
	ManualReview            bool   `json:"manual_review,omitempty"`
	RetainHeader            bool   `json:"retain_header,omitempty"`
	ShowLineEndings         bool   `json:"show_line_endings,omitempty"`
	ShowTabs                bool   `json:"show_tabs,omitempty"`
	ShowWhitespaceErrors    bool   `json:"show_whitespace_errors,omitempty"`
	SkipDeleted             bool   `json:"skip_deleted,omitempty"`
	SkipUncommented         bool   `json:"skip_uncommented,omitempty"`
	SyntaxHighlighting      bool   `json:"syntax_highlighting,omitempty"`
	HideTopMenu             bool   `json:"hide_top_menu,omitempty"`
	AutoHideDiffTableHeader bool   `json:"auto_hide_diff_table_header,omitempty"`
	HideLineNumbers         bool   `json:"hide_line_numbers,omitempty"`
	TabSize                 int    `json:"tab_size"`
	HideEmptyPane           bool   `json:"hide_empty_pane,omitempty"`
}

type EditPreferencesInfo struct {
	TabSize              int  `json:"tab_size"`
	LineLength           int  `json:"line_length"`
	IndentUnit           int  `json:"indent_unit"`
	CursorBlinkRate      int  `json:"cursor_blink_rate"`
	HideTopMenu          bool `json:"hide_top_menu"`
	ShowTabs             bool `json:"show_tabs"`
	ShowWhitespaceErrors bool `json:"show_whitespace_errors"`
	SyntaxHighlighting   bool `json:"syntax_highlighting"`
	HideLineNumbers      bool `json:"hide_line_numbers"`
	MatchBrackets        bool `json:"match_brackets"`
	LineWrapping         bool `json:"line_wrapping"`
	IndentWithTabs       bool `json:"indent_with_tabs"`
	AutoCloseBrackets    bool `json:"auto_close_brackets"`
	ShowBase             bool `json:"show_base"`
}

type EditPreferencesInput struct {
	TabSize              int  `json:"tab_size"`
	LineLength           int  `json:"line_length"`
	IndentUnit           int  `json:"indent_unit"`
	CursorBlinkRate      int  `json:"cursor_blink_rate"`
	HideTopMenu          bool `json:"hide_top_menu,omitempty"`
	ShowTabs             bool `json:"show_tabs,omitempty"`
	ShowWhitespaceErrors bool `json:"show_whitespace_errors,omitempty"`
	SyntaxHighlighting   bool `json:"syntax_highlighting,omitempty"`
	HideLineNumbers      bool `json:"hide_line_numbers,omitempty"`
	MatchBrackets        bool `json:"match_brackets,omitempty"`
	LineWrapping         bool `json:"line_wrapping,omitempty"`
	IndentWithTabs       bool `json:"indent_with_tabs,omitempty"`
	AutoCloseBrackets    bool `json:"auto_close_brackets,omitempty"`
	ShowBase             bool `json:"show_base,omitempty"`
}

// DiffPreferencesInput entity contains information for setting the diff preferences of a user.
// Fields which are not set will not be updated.
type DiffPreferencesInput struct {
	Context                 int    `json:"context,omitempty"`
	ExpandAllComments       bool   `json:"expand_all_comments,omitempty"`
	IgnoreWhitespace        string `json:"ignore_whitespace,omitempty"`
	IntralineDifference     bool   `json:"intraline_difference,omitempty"`
	LineLength              int    `json:"line_length,omitempty"`
	ManualReview            bool   `json:"manual_review,omitempty"`
	RetainHeader            bool   `json:"retain_header,omitempty"`
	ShowLineEndings         bool   `json:"show_line_endings,omitempty"`
	ShowTabs                bool   `json:"show_tabs,omitempty"`
	ShowWhitespaceErrors    bool   `json:"show_whitespace_errors,omitempty"`
	SkipDeleted             bool   `json:"skip_deleted,omitempty"`
	SkipUncommented         bool   `json:"skip_uncommented,omitempty"`
	SyntaxHighlighting      bool   `json:"syntax_highlighting,omitempty"`
	HideTopMenu             bool   `json:"hide_top_menu,omitempty"`
	AutoHideDiffTableHeader bool   `json:"auto_hide_diff_table_header,omitempty"`
	HideLineNumbers         bool   `json:"hide_line_numbers,omitempty"`
	TabSize                 int    `json:"tab_size,omitempty"`
}

// PreferencesInfo entity contains information about a user’s preferences.
type PreferencesInfo struct {
	ChangesPerPage            int               `json:"changes_per_page"`
	ShowSiteHeader            bool              `json:"show_site_header,omitempty"`
	UseFlashClipboard         bool              `json:"use_flash_clipboard,omitempty"`
	DownloadScheme            string            `json:"download_scheme"`
	DownloadCommand           string            `json:"download_command"`
	CopySelfOnEmail           bool              `json:"copy_self_on_email,omitempty"`
	DateFormat                string            `json:"date_format"`
	TimeFormat                string            `json:"time_format"`
	RelativeDateInChangeTable bool              `json:"relative_date_in_change_table,omitempty"`
	SizeBarInChangeTable      bool              `json:"size_bar_in_change_table,omitempty"`
	LegacycidInChangeTable    bool              `json:"legacycid_in_change_table,omitempty"`
	MuteCommonPathPrefixes    bool              `json:"mute_common_path_prefixes,omitempty"`
	ReviewCategoryStrategy    string            `json:"review_category_strategy"`
	DiffView                  string            `json:"diff_view"`
	My                        []TopMenuItemInfo `json:"my"`
	URLAliases                string            `json:"url_aliases,omitempty"`
}

// PreferencesInput entity contains information for setting the user preferences.
// Fields which are not set will not be updated.
type PreferencesInput struct {
	ChangesPerPage            int               `json:"changes_per_page,omitempty"`
	ShowSiteHeader            bool              `json:"show_site_header,omitempty"`
	UseFlashClipboard         bool              `json:"use_flash_clipboard,omitempty"`
	DownloadScheme            string            `json:"download_scheme,omitempty"`
	DownloadCommand           string            `json:"download_command,omitempty"`
	CopySelfOnEmail           bool              `json:"copy_self_on_email,omitempty"`
	DateFormat                string            `json:"date_format,omitempty"`
	TimeFormat                string            `json:"time_format,omitempty"`
	RelativeDateInChangeTable bool              `json:"relative_date_in_change_table,omitempty"`
	SizeBarInChangeTable      bool              `json:"size_bar_in_change_table,omitempty"`
	LegacycidInChangeTable    bool              `json:"legacycid_in_change_table,omitempty"`
	MuteCommonPathPrefixes    bool              `json:"mute_common_path_prefixes,omitempty"`
	ReviewCategoryStrategy    string            `json:"review_category_strategy,omitempty"`
	DiffView                  string            `json:"diff_view,omitempty"`
	My                        []TopMenuItemInfo `json:"my,omitempty"`
	URLAliases                string            `json:"url_aliases,omitempty"`
}

// CapabilityOptions specifies the parameters to filter for capabilities.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#list-account-capabilities
type CapabilityOptions struct {
	// To filter the set of global capabilities the q parameter can be used.
	// Filtering may decrease the response time by avoiding looking at every possible alternative for the caller.
	Filter []string `url:"q,omitempty"`
}

// Query lists accounts visible to the caller.
// The query string must be provided by the q parameter.
// The n parameter can be used to limit the returned results.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#query-accounts
func (s *AccountsService) Query(ctx context.Context, opt *QueryAccountOptions) (*[]AccountInfo, *http.Response, error) {
	v := new([]AccountInfo)
	resp, err := s.gerrit.Requester.Call(ctx, "GET", "accounts/", opt, v)
	return v, resp, err
}

// Get returns an account as an AccountInfo entity.
// If account is "self" the current authenticated account will be returned.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#get-account
func (s *AccountsService) Get(ctx context.Context, accountID string) (*Account, *http.Response, error) {
	account := &Account{Raw: new(AccountInfo), gerrit: s.gerrit, Base: accountID}

	resp, err := account.Poll(ctx)
	if err != nil {
		return nil, resp, err
	}

	//if accountID == "self" {
	//        account.Base = strconv.Itoa(account.Raw.AccountID)
	//}

	return account, resp, nil
}

// Create creates a new account.
// In the request body additional data for the account can be provided as AccountInput.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#create-account
func (s *AccountsService) Create(ctx context.Context, Username string, input *AccountInput) (*Account, *http.Response, error) {
	obj := Account{Raw: new(AccountInfo), gerrit: s.gerrit, Base: Username}
	return obj.Create(ctx, input)
}

func (a *Account) Poll(ctx context.Context) (*http.Response, error) {
	u := fmt.Sprintf("accounts/%s", a.Base)
	return a.gerrit.Requester.Call(ctx, "GET", u, nil, a.Raw)
}

// Create creates a new account.
// In the request body additional data for the account can be provided as AccountInput.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#create-account
func (a *Account) Create(ctx context.Context, input *AccountInput) (*Account, *http.Response, error) {
	v := new(AccountInfo)
	u := fmt.Sprintf("accounts/%s", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "PUT", u, input, v)
	if err != nil {
		return nil, resp, err
	}

	resp, err = a.Poll(ctx)
	if err != nil {
		return nil, resp, err
	}

	a.Base = strconv.Itoa(v.AccountID)
	return a, resp, nil
}

// GetDetails retrieves the details of an account.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#get-detail
func (a *Account) GetDetails(ctx context.Context) (*AccountDetailInfo, *http.Response, error) {
	v := new(AccountDetailInfo)
	u := fmt.Sprintf("accounts/%s/detail", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetName retrieves the full name of an account.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#get-account-name
func (a *Account) GetName(ctx context.Context) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("accounts/%s/name", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// SetName sets the full name of an account.
// The new account name must be provided in the request body inside an AccountNameInput entity.
//
// As response the new account name is returned.
// If the name was deleted the response is “204 No Content”.
// Some realms may not allow to modify the account name.
// In this case the request is rejected with “405 Method Not Allowed”.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#set-account-name
func (a *Account) SetName(ctx context.Context, input *AccountNameInput) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("accounts/%s/name", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "PUT", u, input, v)
	if err != nil {
		return "", resp, err
	}

	return *v, resp, nil
}

// DeleteName deletes the name of an account.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#delete-account-name
func (a *Account) DeleteName(ctx context.Context) (*http.Response, error) {
	u := fmt.Sprintf("accounts/%s/name", a.Base)
	return a.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)
}

// GetStatus Retrieves the status of an account.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#get-account-status
func (a *Account) GetStatus(ctx context.Context) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("accounts/%s/status", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// SetStatus Sets the status of an account.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#set-account-status
func (a *Account) SetStatus(ctx context.Context, input *AccountStatusInput) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("accounts/%s/status", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "PUT", u, input, v)
	if err != nil {
		return "", resp, err
	}

	return *v, resp, nil
}

// GetUsername retrieves the username of an account.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#get-username
func (a *Account) GetUsername(ctx context.Context) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("accounts/%s/username", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// SetUsername sets a new username.
// The new username must be provided in the request body inside a UsernameInput entity.
// Once set, the username cannot be changed or deleted.
// If attempted this fails with “405 Method Not Allowed”.
//
// As response the new username is returned.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#set-username
func (a *Account) SetUsername(ctx context.Context, input *UsernameInput) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("accounts/%s/username", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "PUT", u, input, v)
	if err != nil {
		return "", resp, err
	}

	return *v, resp, nil
}

func (a *Account) SetDisplayName(ctx context.Context, input *DisplayNameInput) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("accounts/%s/displayname", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "PUT", u, input, v)
	if err != nil {
		return "", resp, err
	}

	return *v, resp, nil
}

// GetActive checks if an account is active.
//
// If the account is active the string ok is returned.
// If the account is inactive the response is “204 No Content”.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#get-active
func (a *Account) GetActive(ctx context.Context) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("accounts/%s/active", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// SetActive sets the account state to active.
//
// If the account was already active the response is “200 OK”.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#set-active
func (a *Account) SetActive(ctx context.Context) (*http.Response, error) {
	u := fmt.Sprintf("accounts/%s/active", a.Base)
	return a.gerrit.Requester.Call(ctx, "PUT", u, nil, nil)
}

// DeleteActive sets the account state to inactive.
// If the account was already inactive the response is “404 Not Found”.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#delete-active
func (a *Account) DeleteActive(ctx context.Context) (*http.Response, error) {
	u := fmt.Sprintf("accounts/%s/active", a.Base)
	return a.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)
}

// GetHTTPPassword retrieves the HTTP password of an account.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#get-http-password
func (a *Account) GetHTTPPassword(ctx context.Context) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("accounts/%s/password.http", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// SetHTTPPassword sets/Generates the HTTP password of an account.
// The options for setting/generating the HTTP password must be provided in the request body inside a HTTPPasswordInput entity.
//
// As response the new HTTP password is returned.
// If the HTTP password was deleted the response is “204 No Content”.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#set-http-password
func (a *Account) SetHTTPPassword(ctx context.Context, input *HTTPPasswordInput) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("accounts/%s/password.http", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "PUT", u, input, v)
	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// DeleteHTTPPassword deletes the HTTP password of an account.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#delete-http-password
func (a *Account) DeleteHTTPPassword(ctx context.Context) (*http.Response, error) {
	u := fmt.Sprintf("accounts/%s/password.http", a.Base)
	return a.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)
}

// GetOAuthAccessToken Returns a previously obtained OAuth access token.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#get-oauth-token
func (a *Account) GetOAuthAccessToken(ctx context.Context) (*OAuthTokenInfo, *http.Response, error) {
	v := new(OAuthTokenInfo)
	u := fmt.Sprintf("accounts/%s/oauthtoken", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// ListEmails returns the email addresses that are configured for the specified user.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#list-account-emails
func (a *Account) ListEmails(ctx context.Context) (*[]EmailInfo, *http.Response, error) {
	v := new([]EmailInfo)
	u := fmt.Sprintf("accounts/%s/emails", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetEmail retrieves an email address of a user.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#get-account-email
func (a *Account) GetEmail(ctx context.Context, emailID string) (*EmailInfo, *http.Response, error) {
	v := new(EmailInfo)
	u := fmt.Sprintf("accounts/%s/emails/%s", a.Base, emailID)

	resp, err := a.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// CreateEmail registers a new email address for the user.
// A verification email is sent with a link that needs to be visited to confirm the email address, unless DEVELOPMENT_BECOME_ANY_ACCOUNT is used as authentication type.
// For the development mode email addresses are directly added without confirmation.
// A Gerrit administrator may add an email address without confirmation by setting no_confirmation in the EmailInput.
// In the request body additional data for the email address can be provided as EmailInput.
//
// As response the new email address is returned as EmailInfo entity.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#create-account-email
func (a *Account) CreateEmail(ctx context.Context, emailID string, input *EmailInput) (*EmailInfo, *http.Response, error) {
	v := new(EmailInfo)
	u := fmt.Sprintf("accounts/%s/emails/%s", a.Base, emailID)

	resp, err := a.gerrit.Requester.Call(ctx, "PUT", u, input, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// DeleteEmail deletes an email address of an account.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#delete-account-email
func (a *Account) DeleteEmail(ctx context.Context, emailID string) (*http.Response, error) {
	u := fmt.Sprintf("accounts/%s/emails/%s", a.Base, emailID)
	return a.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)
}

// SetPreferredEmail sets an email address as preferred email address for an account.
//
// If the email address was already the preferred email address of the account the response is “200 OK”.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#set-preferred-email
func (a *Account) SetPreferredEmail(ctx context.Context, emailID string) (*http.Response, error) {
	u := fmt.Sprintf("accounts/%s/emails/%s/preferred", a.Base, emailID)
	return a.gerrit.Requester.Call(ctx, "PUT", u, nil, nil)
}

// ListSSHKeys returns the SSH keys of an account.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#list-ssh-keys
func (a *Account) ListSSHKeys(ctx context.Context) (*[]SSHKeyInfo, *http.Response, error) {
	v := new([]SSHKeyInfo)
	u := fmt.Sprintf("accounts/%s/sshkeys", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetSSHKey retrieves an SSH key of a user.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#get-ssh-key
func (a *Account) GetSSHKey(ctx context.Context, sshKeyID string) (*SSHKeyInfo, *http.Response, error) {
	v := new(SSHKeyInfo)
	u := fmt.Sprintf("accounts/%s/sshkeys/%s", a.Base, sshKeyID)

	resp, err := a.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// AddSSHKey adds an SSH key for a user.
// The SSH public key must be provided as raw content in the request body.
// Trying to add an SSH key that already exists succeeds, but no new SSH key is persisted.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#add-ssh-key
func (a *Account) AddSSHKey(ctx context.Context, sshKey string) (*SSHKeyInfo, *http.Response, error) {
	v := new(SSHKeyInfo)
	u := fmt.Sprintf("accounts/%s/sshkeys", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "POST", u, sshKey, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// DeleteSSHKey deletes an SSH key of a user.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#delete-ssh-key
func (a *Account) DeleteSSHKey(ctx context.Context, sshKeyID int) (*http.Response, error) {
	u := fmt.Sprintf("accounts/%s/sshkeys/%d", a.Base, sshKeyID)
	return a.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)
}

// ListGPGKeys returns the GPG keys of an account.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#list-gpg-keys
func (a *Account) ListGPGKeys(ctx context.Context) (*map[string]GpgKeyInfo, *http.Response, error) {
	v := new(map[string]GpgKeyInfo)
	u := fmt.Sprintf("accounts/%s/gpgkeys", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// AddGPGKey Add or delete one or more GPG keys for a user.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#add-gpg-key
func (a *Account) AddGPGKey(ctx context.Context, input *GpgKeysInput) (map[string]GpgKeyInfo, *http.Response, error) {
	v := make(map[string]GpgKeyInfo)
	u := fmt.Sprintf("accounts/%s/gpgkeys", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "POST", u, input, &v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetGPGKey retrieves a GPG key of a user.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#get-gpg-key
func (a *Account) GetGPGKey(ctx context.Context, gpgKeyID string) (*GpgKeyInfo, *http.Response, error) {
	v := new(GpgKeyInfo)
	u := fmt.Sprintf("accounts/%s/gpgkeys/%s", a.Base, gpgKeyID)

	resp, err := a.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// DeleteGPGKey deletes a GPG key of a user.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#delete-gpg-key
func (a *Account) DeleteGPGKey(ctx context.Context, gpgKeyID string) (*http.Response, error) {
	u := fmt.Sprintf("accounts/%s/gpgkeys/%s", a.Base, gpgKeyID)
	return a.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)
}

// ListCapabilities returns the global capabilities that are enabled for the specified user.
// If the global capabilities for the calling user should be listed, self can be used as account-id.
// This can be used by UI tools to discover if administrative features are available to the caller, so they can hide (or show) relevant UI actions.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#list-account-capabilities
func (a *Account) ListCapabilities(ctx context.Context, opt *CapabilityOptions) (*AccountCapabilityInfo, *http.Response, error) {
	v := new(AccountCapabilityInfo)
	u := fmt.Sprintf("accounts/%s/capabilities", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "GET", u, opt, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// CheckCapability checks if a user has a certain global capability.
//
// If the user has the global capability the string ok is returned.
// If the user doesn’t have the global capability the response is “404 Not Found”.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#check-account-capability
func (a *Account) CheckCapability(ctx context.Context, capabilityID string) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("accounts/%s/capabilities/%s", a.Base, capabilityID)

	resp, err := a.gerrit.Requester.Call(ctx, "GET", u, nil, v)

	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// ListGroups lists all groups that contain the specified user as a member.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#list-groups
func (a *Account) ListGroups(ctx context.Context) (*[]GroupInfo, *http.Response, error) {
	v := new([]GroupInfo)
	u := fmt.Sprintf("accounts/%s/groups/", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetAvatar

// GetAvatarChangeURL retrieves the URL where the user can change the avatar image.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#get-avatar-change-url
func (a *Account) GetAvatarChangeURL(ctx context.Context) (string, *http.Response, error) {
	v := new(string)
	u := fmt.Sprintf("accounts/%s/avatar.change.url", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return "", resp, err
	}
	return *v, resp, nil
}

// GetUserPreferences retrieves the user’s preferences.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#get-user-preferences
func (a *Account) GetUserPreferences(ctx context.Context) (*PreferencesInfo, *http.Response, error) {
	v := new(PreferencesInfo)
	u := fmt.Sprintf("accounts/%s/preferences", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// SetUserPreferences sets the user’s preferences.
// The new preferences must be provided in the request body as a PreferencesInput entity.
//
// As result the new preferences of the user are returned as a PreferencesInfo entity.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#set-user-preferences
func (a *Account) SetUserPreferences(ctx context.Context, input *PreferencesInput) (*PreferencesInfo, *http.Response, error) {
	v := new(PreferencesInfo)
	u := fmt.Sprintf("accounts/%s/preferences", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "PUT", u, input, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetDiffPreferences retrieves the diff preferences of a user.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#get-diff-preferences
func (a *Account) GetDiffPreferences(ctx context.Context) (*DiffPreferencesInfo, *http.Response, error) {
	v := new(DiffPreferencesInfo)
	u := fmt.Sprintf("accounts/%s/preferences.diff", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// SetDiffPreferences sets the diff preferences of a user.
// The new diff preferences must be provided in the request body as a DiffPreferencesInput entity.
//
// As result the new diff preferences of the user are returned as a DiffPreferencesInfo entity.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#set-diff-preferences
func (a *Account) SetDiffPreferences(ctx context.Context, input *DiffPreferencesInput) (*DiffPreferencesInfo, *http.Response, error) {
	v := new(DiffPreferencesInfo)
	u := fmt.Sprintf("accounts/%s/preferences.diff", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "PUT", u, input, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetEditPreferences retrieves the edit preferences of a user.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#get-edit-preferences
func (a *Account) GetEditPreferences(ctx context.Context) (*EditPreferencesInfo, *http.Response, error) {
	v := new(EditPreferencesInfo)
	u := fmt.Sprintf("accounts/%s/preferences.edit", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// SetEditPreferences sets the edit preferences of a user.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#set-edit-preferences
func (a *Account) SetEditPreferences(ctx context.Context, input *EditPreferencesInput) (*EditPreferencesInfo, *http.Response, error) {
	v := new(EditPreferencesInfo)
	u := fmt.Sprintf("accounts/%s/preferences.edit", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "PUT", u, input, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetExternalIDs retrieves the external ids of a user account.
//
// Only external ids belonging to the caller may be requested.
// Users that have Modify Account can request external ids that belong to other accounts.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#get-account-external-ids
func (a *Account) GetExternalIDs(ctx context.Context) (*[]AccountExternalIdInfo, *http.Response, error) {
	v := new([]AccountExternalIdInfo)
	u := fmt.Sprintf("accounts/%s/external.ids", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// GetStarredChanges gets the changes starred by the identified user account.
// This URL endpoint is functionally identical to the changes query GET /changes/?q=is:starred.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#get-starred-changes
func (a *Account) GetStarredChanges(ctx context.Context) (*[]ChangeInfo, *http.Response, error) {
	v := new([]ChangeInfo)
	u := fmt.Sprintf("accounts/%s/starred.changes", a.Base)

	resp, err := a.gerrit.Requester.Call(ctx, "GET", u, nil, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

// StarChange star a change.
// Starred changes are returned for the search query is:starred or starredby:USER and automatically notify the user whenever updates are made to the change.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#star-change
func (a *Account) StarChange(ctx context.Context, changeID string) (*http.Response, error) {
	u := fmt.Sprintf("accounts/%s/starred.changes/%s", a.Base, changeID)
	return a.gerrit.Requester.Call(ctx, "PUT", u, nil, nil)
}

// UnstarChange nstar a change.
// Removes the starred flag, stopping notifications.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html#unstar-change
func (a *Account) UnstarChange(ctx context.Context, changeID string) (*http.Response, error) {
	u := fmt.Sprintf("accounts/%s/starred.changes/%s", a.Base, changeID)
	return a.gerrit.Requester.Call(ctx, "DELETE", u, nil, nil)
}