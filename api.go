package confluence

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2021 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Content type
const (
	CONTENT_TYPE_ATTACHMENT = "attachment"
	CONTENT_TYPE_BLOGPOST   = "blogpost"
	CONTENT_TYPE_COMMENT    = "comment"
	CONTENT_TYPE_PAGE       = "page"
)

// Excerpt values
const (
	SEARCH_EXCERPT_INDEXED   = "indexed"
	SEARCH_EXCERPT_HIGHLIGHT = "highlight"
	SEARCH_EXCERPT_NONE      = "none"
)

// Space type
const (
	SPACE_TYPE_PERSONAL = "personal"
	SPACE_TYPE_GLOBAL   = "global"
)

// Content status
const (
	SPACE_STATUS_CURRENT  = "current"
	SPACE_STATUS_ARCHIVED = "archived"
)

// Space status
const (
	CONTENT_STATUS_CURRENT = "current"
	CONTENT_STATUS_TRASHED = "trashed"
	CONTENT_STATUS_DRAFT   = "draft"
)

// Units
const (
	UNITS_MINUTES = "minutes"
	UNITS_HOURS   = "hours"
	UNITS_DAYS    = "days"
	UNITS_MONTHS  = "months"
	UNITS_YEARS   = "years"
)

// Operations types
const (
	OPERATION_READ   = "read"
	OPERATION_UPDATE = "update"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Parameters is interface for parameters structs
type Parameters interface {
	ToQuery() string
	Validate() error
}

// Date is RFC3339 encoded date
type Date struct {
	time.Time
}

// Timestamp is UNIX timestamp in ms
type Timestamp struct {
	time.Time
}

// ContainerID is container ID
type ContainerID string

// ExtensionPosition is extension position
type ExtensionPosition int

// EmptyParameters is empty parameters
type EmptyParameters struct {
	// nothing
}

// ExpandParameters is params with field expand info
type ExpandParameters struct {
	Expand []string `query:"expand"`
}

// CollectionParameters is params with pagination info
type CollectionParameters struct {
	Expand []string `query:"expand"`
	Start  int      `query:"start"`
	Limit  int      `query:"limit"`
}

// AUDIT ///////////////////////////////////////////////////////////////////////////////

// AuditParameters is params for fetching audit data
type AuditParameters struct {
	StartDate    time.Time `query:"startDate"`
	EndDate      time.Time `query:"endDate"`
	SearchString string    `query:"searchString"`
	Start        int       `query:"start"`
	Limit        int       `query:"limit"`
}

// AuditSinceParameters is params for fetching audit data
type AuditSinceParameters struct {
	Number       int    `query:"number"`
	Units        string `query:"units"`
	SearchString string `query:"searchString"`
	Start        int    `query:"start"`
	Limit        int    `query:"limit"`
}

// AuditRecord represents audit record
type AuditRecord struct {
	Author        *User      `json:"author"`
	RemoteAddress string     `json:"remoteAddress"`
	CreationDate  *Timestamp `json:"creationDate"`
	Summary       string     `json:"summary"`
	Description   string     `json:"description"`
	Category      string     `json:"category"`
	IsSysAdmin    bool       `json:"sysAdmin"`
}

// AuditRecordCollection contains paginated list of audit record
type AuditRecordCollection struct {
	Results []*AuditRecord `json:"results"`
	Start   int            `json:"start"`
	Limit   int            `json:"limit"`
	Size    int            `json:"size"`
}

// AuditRetentionInfo contains info about retention time
type AuditRetentionInfo struct {
	Number int    `json:"number"`
	Units  string `json:"units"`
}

// ATTACHMENTS /////////////////////////////////////////////////////////////////////////

// AttachmentParameters is params for fetching attachments info
type AttachmentParameters struct {
	Filename  string   `query:"filename"`
	MediaType string   `query:"mediaType"`
	Expand    []string `query:"expand"`
	Start     int      `query:"start"`
	Limit     int      `query:"limit"`
}

// CONTENT /////////////////////////////////////////////////////////////////////////////

// ContentParameters is params for fetching content info
type ContentParameters struct {
	Type       string    `query:"type"`
	SpaceKey   string    `query:"spaceKey"`
	Title      string    `query:"title"`
	Status     string    `query:"status"`
	PostingDay time.Time `query:"postingDay"`
	Expand     []string  `query:"expand"`
	Start      int       `query:"start"`
	Limit      int       `query:"limit"`
}

// ContentIDParameters is params for fetching content info
type ContentIDParameters struct {
	Status  string   `query:"status"`
	Version int      `query:"version"`
	Expand  []string `query:"expand"`
}

// ContentSearchParameters is params for searching content
type ContentSearchParameters struct {
	CQL        string   `query:"cql"`
	CQLContext string   `query:"cqlcontext"`
	Expand     []string `query:"expand"`
	Start      int      `query:"start"`
	Limit      int      `query:"limit"`
}

// ChildrenParameters is params for fetching content child info
type ChildrenParameters struct {
	ParentVersion int      `query:"parentVersion"`
	Location      string   `query:"location"`
	Depth         string   `query:"depth"`
	Expand        []string `query:"expand"`
	Start         int      `query:"start"`
	Limit         int      `query:"limit"`
}

// Content contains content info
type Content struct {
	ID          string       `json:"id"`
	Type        string       `json:"type"`
	Status      string       `json:"status"`
	Title       string       `json:"title"`
	Extensions  *Extensions  `json:"extensions"`
	Metadata    *Metadata    `json:"metadata"`
	Container   *Container   `json:"container"`
	Space       *Space       `json:"space"`
	Version     *Version     `json:"version"`
	Operations  []*Operation `json:"operations"`
	Children    *Contents    `json:"children"`
	Ancestors   []*Content   `json:"ancestors"`
	Descendants *Contents    `json:"descendants"`
	Body        *Body        `json:"body"`
	Links       *Links       `json:"_links"`
}

// ContentCollection represents paginated list of content
type ContentCollection struct {
	Results []*Content `json:"results"`
	Start   int        `json:"start"`
	Limit   int        `json:"limit"`
	Size    int        `json:"size"`
}

// Contents contains all types of content
type Contents struct {
	Attachments *ContentCollection `json:"attachment"`
	Comments    *ContentCollection `json:"comment"`
	Pages       *ContentCollection `json:"page"`
	Blogposts   *ContentCollection `json:"blogposts"`
}

// Body contains content data
type Body struct {
	View        *View `json:"view"`
	ExportView  *View `json:"export_view"`
	StyledView  *View `json:"styled_view"`
	StorageView *View `json:"storage"`
}

// View is data view
type View struct {
	Representation string `json:"representation"`
	Value          string `json:"value"`
}

// Version contains info about content version
type Version struct {
	Message     string   `json:"message"`
	By          *User    `json:"by"`
	When        *Date    `json:"when"`
	Number      int      `json:"number"`
	Content     *Content `json:"content"`
	IsMinorEdit bool     `json:"minorEdit"`
	IsHidden    bool     `json:"hidden"`
}

// Extensions contains info about content extensions
type Extensions struct {
	Position   ExtensionPosition `json:"position"`   // Page
	MediaType  string            `json:"mediaType"`  // Attachment
	FileSize   int               `json:"fileSize"`   // Attachment
	Comment    string            `json:"comment"`    // Attachment
	Location   string            `json:"location"`   // Comment
	Resolution *Resolution       `json:"resolution"` // Comment
}

// Resolution contains resolution info
type Resolution struct {
	Status           string `json:"status"`
	LastModifier     *User  `json:"lastModifier"`
	LastModifiedDate *Date  `json:"lastModifiedDate"`
}

// Operation contains operation info
type Operation struct {
	Name       string `json:"operation"`
	TargetType string `json:"targetType"`
}

// Metadata contains metadata records
type Metadata struct {
	Labels    *LabelCollection `json:"labels"`    // Page
	MediaType string           `json:"mediaType"` // Attachment
}

// History contains info about content history
type History struct {
	CreatedBy       *User         `json:"createdBy"`
	CreatedDate     *Date         `json:"createdDate"`
	LastUpdated     *Version      `json:"lastUpdated"`
	PreviousVersion *Version      `json:"previousVersion"`
	NextVersion     *Version      `json:"nextVersion"`
	Contributors    *Contributors `json:"contributors"`
	IsLatest        bool          `json:"latest"`
}

// Contributors contains contributors list
type Contributors struct {
	Publishers *Publishers `json:"publishers"`
}

// Publishers contains info about users
type Publishers struct {
	Users    []*User  `json:"users"`
	UserKeys []string `json:"userKeys"`
}

// Container contains basic container info
type Container struct {
	ID    ContainerID `json:"id"`
	Key   string      `json:"key"`   // Space
	Name  string      `json:"name"`  // Space
	Title string      `json:"title"` // Page or blogpost
	Links *Links      `json:"_links"`
}

// LABELS //////////////////////////////////////////////////////////////////////////////

// LabelParameters is params for fetching labels
type LabelParameters struct {
	Prefix string `query:"prefix"`
	Start  int    `query:"start"`
	Limit  int    `query:"limit"`
}

// LabelCollection contains paginated list of labels
type LabelCollection struct {
	Result []*Label `json:"results"`
	Start  int      `json:"start"`
	Limit  int      `json:"limit"`
	Size   int      `json:"size"`
}

// Label contains label info
type Label struct {
	Prefix string `json:"prefix"`
	Name   string `json:"name"`
	ID     string `json:"id"`
}

// GROUPS //////////////////////////////////////////////////////////////////////////////

// Group contains group info
type Group struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// GroupCollection contains paginated list of groups
type GroupCollection struct {
	Results []*Group `json:"results"`
	Start   int      `json:"start"`
	Limit   int      `json:"limit"`
	Size    int      `json:"size"`
}

// RESTRICTIONS ////////////////////////////////////////////////////////////////////////

// Restrictions contains info about all restrictions
type Restrictions struct {
	Read   *Restriction `json:"read"`
	Update *Restriction `json:"update"`
}

// Restriction contains restriction info for single operation
type Restriction struct {
	Operation string           `json:"operation"`
	Data      *RestrictionData `json:"restrictions"`
}

// RestrictionData contains restrictions data
type RestrictionData struct {
	User  *UserCollection  `json:"user"`
	Group *GroupCollection `json:"group"`
}

// SEARCH //////////////////////////////////////////////////////////////////////////////

// SearchParameters is params for fetching search results
type SearchParameters struct {
	Expand                []string `query:"expand"`
	CQL                   string   `query:"cql"`
	CQLContext            string   `query:"cqlcontext"`
	Excerpt               string   `query:"excerpt"`
	Start                 int      `query:"start"`
	Limit                 int      `query:"limit"`
	IncludeArchivedSpaces bool     `query:"includeArchivedSpaces"`
}

// SearchResult contains contains paginated list of search results
type SearchResult struct {
	Results        []*SearchEntity `json:"results"`
	Start          int             `json:"start"`
	Limit          int             `json:"limit"`
	Size           int             `json:"size"`
	TotalSize      int             `json:"totalSize"`
	CQLQuery       string          `json:"cqlQuery"`
	SearchDuration int             `json:"searchDuration"`
}

// SearchEntity contains search result
type SearchEntity struct {
	Content      *Content `json:"content"`
	Space        *Space   `json:"space"`
	User         *User    `json:"user"`
	Title        string   `json:"title"`
	Excerpt      string   `json:"excerpt"`
	URL          string   `json:"url"`
	EntityType   string   `json:"entityType"`
	LastModified *Date    `json:"lastModified"`
}

// SPACE ///////////////////////////////////////////////////////////////////////////////

// SpaceParameters is params for fetching info about space
type SpaceParameters struct {
	SpaceKey  []string `query:"spaceKey,unwrap"`
	Expand    []string `query:"expand"`
	Type      string   `query:"type"`
	Status    string   `query:"status"`
	Label     string   `query:"label"`
	Depth     string   `query:"depth"`
	Start     int      `query:"start"`
	Limit     int      `query:"limit"`
	Favourite bool     `query:"favourite"`
}

// Space contains info about space
type Space struct {
	ID    int    `json:"id"`
	Key   string `json:"key"`
	Name  string `json:"name"`
	Icon  *Icon  `json:"icon"`
	Type  string `json:"type"`
	Links *Links `json:"_links"`
}

// SpaceCollection contains paginated list of spaces
type SpaceCollection struct {
	Results []*Space `json:"results"`
	Start   int      `json:"start"`
	Limit   int      `json:"limit"`
	Size    int      `json:"size"`
}

// Icon contains icon info
type Icon struct {
	Path      string `json:"path"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	IsDefault bool   `json:"isDefault"`
}

// USER ////////////////////////////////////////////////////////////////////////////////

// UserParameters is params for fetching info about user
type UserParameters struct {
	Key      string   `query:"key"`
	Username string   `query:"username"`
	Expand   []string `query:"expand"`
	Start    int      `query:"start"`
	Limit    int      `query:"limit"`
}

// User contains user info
type User struct {
	Type           string `json:"type"`
	Name           string `json:"username"`
	Key            string `json:"userKey"`
	ProfilePicture *Icon  `json:"profilePicture"`
	DisplayName    string `json:"displayName"`
}

// UserCollection contains paginated list of users
type UserCollection struct {
	Results []*User `json:"results"`
	Start   int     `json:"start"`
	Limit   int     `json:"limit"`
	Size    int     `json:"size"`
}

// LINKS ///////////////////////////////////////////////////////////////////////////////

// Links contains links
type Links struct {
	WebUI  string `json:"webui"`
	TinyUI string `json:"tinyui"`
	Base   string `json:"base"`
}

// WATCH ///////////////////////////////////////////////////////////////////////////////

// WatchParameters is params for fetching info about watchers
type WatchParameters struct {
	Key         string `query:"key"`
	Username    string `query:"username"`
	ContentType string `query:"contentType"`
}

// ListWatchersParameters is params for fetching info about page watchers
type ListWatchersParameters struct {
	PageID string `query:"pageId"`
}

// WatchStatus contains watching status
type WatchStatus struct {
	IsWatching bool `json:"watching"`
}

// WatchInfo contains info about watchers
type WatchInfo struct {
	PageWatchers  []*Watcher `json:"pageWatchers"`
	SpaceWatchers []*Watcher `json:"spaceWatchers"`
}

// Watcher contains watcher info
type Watcher struct {
	AvatarURL   string `json:"avatarUrl"`
	Name        string `json:"name"`
	Key         string `json:"userKey"`
	DisplayName string `json:"fullName"`
	Type        string `json:"type"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// IsAttachment return true if content is attachment
func (c *Content) IsAttachment() bool {
	return c.Type == CONTENT_TYPE_ATTACHMENT
}

// IsComment return true if content is comment
func (c *Content) IsComment() bool {
	return c.Type == CONTENT_TYPE_COMMENT
}

// IsPage return true if content is page
func (c *Content) IsPage() bool {
	return c.Type == CONTENT_TYPE_PAGE
}

// IsTrashed return true if content is trashed
func (c *Content) IsTrashed() bool {
	return c.Status == CONTENT_STATUS_TRASHED
}

// IsDraft return true if content is draft
func (c *Content) IsDraft() bool {
	return c.Status == CONTENT_STATUS_DRAFT
}

// IsGlobal return true if space is global
func (s *Space) IsGlobal() bool {
	return s.Type == SPACE_TYPE_GLOBAL
}

// IsPersonal return true if space is personal
func (s *Space) IsPersonal() bool {
	return s.Type == SPACE_TYPE_PERSONAL
}

// IsArchived return true if space is archived
func (s *Space) IsArchived() bool {
	return s.Type == SPACE_STATUS_ARCHIVED
}

// IsPage return true if container is page
func (c *Container) IsPage() bool {
	return c.Title != ""
}

// IsSpace return true if container is space
func (c *Container) IsSpace() bool {
	return c.Key != ""
}

// Combined return united slice with all watchers
func (wi *WatchInfo) Combined() []*Watcher {
	var result []*Watcher

	result = append(result, wi.PageWatchers...)

MAINLOOP:
	for _, watcher := range wi.SpaceWatchers {
		for _, pageWatcher := range wi.PageWatchers {
			if watcher.Key == pageWatcher.Key {
				continue MAINLOOP
			}
		}

		result = append(result, watcher)
	}

	return result
}

// ////////////////////////////////////////////////////////////////////////////////// //

// UnmarshalJSON is custom Date format unmarshaler
func (d *Date) UnmarshalJSON(b []byte) error {
	var err error

	d.Time, err = time.Parse(time.RFC3339, strings.Trim(string(b), "\""))

	if err != nil {
		return fmt.Errorf("Cannot unmarshal Date value: %v", err)
	}

	return nil
}

// UnmarshalJSON is custom container ID unmarshaler
func (c *ContainerID) UnmarshalJSON(b []byte) error {
	switch {
	case len(b) == 0:
		// nop
	case b[0] == '"':
		*c = ContainerID(strings.Replace(string(b), "\"", "", -1))
	default:
		*c = ContainerID(string(b))
	}

	return nil
}

// UnmarshalJSON is custom position unmarshaler
func (ep *ExtensionPosition) UnmarshalJSON(b []byte) error {
	if string(b) == "\"none\"" {
		*ep = ExtensionPosition(-1)
		return nil
	}

	v, err := strconv.Atoi(string(b))

	if err != nil {
		return fmt.Errorf("Cannot unmarshal ExtensionPosition value: %v", err)
	}

	*ep = ExtensionPosition(v)

	return nil
}

// UnmarshalJSON is custom Timestamp format unmarshaler
func (d *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.ParseInt(string(b), 10, 64)

	if err != nil {
		return err
	}

	d.Time = time.Unix(ts/1000, (ts%1000)*1000000)

	if err != nil {
		return fmt.Errorf("Cannot unmarshal Timestamp value: %v", err)
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Validate validates parameters
func (p EmptyParameters) Validate() error {
	return nil
}

// Validate validates parameters
func (p ExpandParameters) Validate() error {
	return nil
}

// Validate validates parameters
func (p CollectionParameters) Validate() error {
	return nil
}

// Validate validates parameters
func (p AuditParameters) Validate() error {
	return nil
}

// Validate validates parameters
func (p AuditSinceParameters) Validate() error {
	return nil
}

// Validate validates parameters
func (p ContentParameters) Validate() error {
	if p.SpaceKey == "" {
		return errors.New("SpaceKey is mandatory and must be set")
	}

	return nil
}

// Validate validates parameters
func (p ContentIDParameters) Validate() error {
	return nil
}

// Validate validates parameters
func (p ContentSearchParameters) Validate() error {
	return nil
}

// Validate validates parameters
func (p ChildrenParameters) Validate() error {
	return nil
}

// Validate validates parameters
func (p AttachmentParameters) Validate() error {
	return nil
}

// Validate validates parameters
func (p LabelParameters) Validate() error {
	return nil
}

// Validate validates parameters
func (p SearchParameters) Validate() error {
	if p.CQL == "" {
		return errors.New("CQL is mandatory and must be set")
	}

	return nil
}

// Validate validates parameters
func (p SpaceParameters) Validate() error {
	if len(p.SpaceKey) == 0 {
		return errors.New("SpaceKey is mandatory and must be set")
	}

	return nil
}

// Validate validates parameters
func (p UserParameters) Validate() error {
	if p.Key == "" && p.Username == "" {
		return errors.New("Key or Username must be set")
	}

	return nil
}

// Validate validates parameters
func (p WatchParameters) Validate() error {
	return nil
}

// Validate validates parameters
func (p ListWatchersParameters) Validate() error {
	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// ToQuery convert params to URL query
func (p EmptyParameters) ToQuery() string {
	return ""
}

// ToQuery convert params to URL query
func (p ExpandParameters) ToQuery() string {
	return paramsToQuery(p)
}

// ToQuery convert params to URL query
func (p CollectionParameters) ToQuery() string {
	return paramsToQuery(p)
}

// ToQuery convert params to URL query
func (p AuditParameters) ToQuery() string {
	return paramsToQuery(p)
}

// ToQuery convert params to URL query
func (p AuditSinceParameters) ToQuery() string {
	return paramsToQuery(p)
}

// ToQuery convert params to URL query
func (p ContentParameters) ToQuery() string {
	return paramsToQuery(p)
}

// ToQuery convert params to URL query
func (p ContentIDParameters) ToQuery() string {
	return paramsToQuery(p)
}

// ToQuery convert params to URL query
func (p ContentSearchParameters) ToQuery() string {
	return paramsToQuery(p)
}

// ToQuery convert params to URL query
func (p ChildrenParameters) ToQuery() string {
	return paramsToQuery(p)
}

// ToQuery convert params to URL query
func (p AttachmentParameters) ToQuery() string {
	return paramsToQuery(p)
}

// ToQuery convert params to URL query
func (p LabelParameters) ToQuery() string {
	return paramsToQuery(p)
}

// ToQuery convert params to URL query
func (p SearchParameters) ToQuery() string {
	return paramsToQuery(p)
}

// ToQuery convert params to URL query
func (p SpaceParameters) ToQuery() string {
	return paramsToQuery(p)
}

// ToQuery convert params to URL query
func (p UserParameters) ToQuery() string {
	return paramsToQuery(p)
}

// ToQuery convert params to URL query
func (p WatchParameters) ToQuery() string {
	return paramsToQuery(p)
}

// ToQuery convert params to URL query
func (p ListWatchersParameters) ToQuery() string {
	return paramsToQuery(p)
}
