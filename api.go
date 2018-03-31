package confluence

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2018 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"strings"
	"time"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	SEARCH_EXCERPT_INDEXED   = "indexed"
	SEARCH_EXCERPT_HIGHLIGHT = "highlight"
	SEARCH_EXCERPT_NONE      = "none"
)

const (
	SPACE_TYPE_PERSONAL = "personal"
	SPACE_TYPE_GLOBAL   = "global"
)

const (
	SPACE_STATUS_CURRENT  = "current"
	SPACE_STATUS_ARCHIVED = "archived"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type Parameters interface {
	ToQuery() string
}

type Date struct {
	time.Time
}

type EmptyParameters struct {
	// nothing
}

type ExpandParameters struct {
	Expand []string `query:"expand"`
}

type CollectionParameters struct {
	Expand []string `query:"expand"`
	Start  int      `query:"start"`
	Limit  int      `query:"limit"`
}

// CONTENT /////////////////////////////////////////////////////////////////////////////

type Content struct {
	ID          string       `json:"id"`
	Type        string       `json:"type"`
	Status      string       `json:"status"`
	Title       string       `json:"title"`
	Extensions  *Extensions  `json:"extensions"`
	Metadata    *Metadata    `json:"metadata"`
	Space       *Space       `json:"space"`
	Version     *Version     `json:"version"`
	Operations  []*Operation `json:"operations"`
	Children    *Contents    `json:"children"`
	Ancestors   []*Content   `json:"ancestors"`
	Descendants *Contents    `json:"descendants"`
	Body        *Body        `json:"body"`
}

type Contents struct {
	Attachments *ContentColletion `json:"attachment"`
	Comments    *ContentColletion `json:"comment"`
	Pages       *ContentColletion `json:"page"`
	Blogposts   *ContentColletion `json:"blogposts"`
}

type Body struct {
	View        *View `json:"view"`
	ExportView  *View `json:"export_view"`
	StyledView  *View `json:"styled_view"`
	StorageView *View `json:"storage"`
}

type View struct {
	Representation string `json:"representation"`
	Value          string `json:"value"`
}

type ContentColletion struct {
	Results []*Content `json:"results"`
	Start   int        `json:"start"`
	Limit   int        `json:"limit"`
	Size    int        `json:"size"`
}

type Version struct {
	By        *User    `json:"by"`
	When      *Date    `json:"when"`
	Message   string   `json:"message"`
	Number    int      `json:"number"`
	MinorEdit bool     `json:"minorEdit"`
	Hidden    bool     `json:"hidden"`
	Content   *Content `json:"content"`
}

type Extensions struct {
	Position   string      `json:"position"`   // Page
	MediaType  string      `json:"mediaType"`  // Attachment
	FileSize   int         `json:"fileSize"`   // Attachment
	Comment    string      `json:"comment"`    // Attachment
	Location   string      `json:"location"`   // Comment
	Resolution *Resolution `json:"resolution"` // Comment
}

type Resolution struct {
	Status           string `json:"status"`
	LastModifier     *User  `json:"lastModifier"`
	LastModifiedDate *Date  `json:"lastModifiedDate"`
}

type Operation struct {
	Name       string `json:"operation"`
	TargetType string `json:"targetType"`
}

type Metadata struct {
	Labels    *LabelCollection `json:"labels"`    // Page
	MediaType string           `json:"mediaType"` // Attachment
}

type LabelCollection struct {
	Result []*Label `json:"results"`
	Start  int      `json:"start"`
	Limit  int      `json:"limit"`
	Size   int      `json:"size"`
}

type Label struct {
	Prefix string `json:"prefix"`
	Name   string `json:"name"`
	ID     string `json:"id"`
}

type History struct {
	Latest          bool          `json:"latest"`
	CreatedBy       *User         `json:"createdBy"`
	CreatedDate     *Date         `json:"createdDate"`
	LastUpdated     *Version      `json:"lastUpdated"`
	PreviousVersion *Version      `json:"previousVersion"`
	NextVersion     *Version      `json:"nextVersion"`
	Contributors    *Contributors `json:"contributors"`
}

type Contributors struct {
	Publishers *Publishers `json:"publishers"`
}

type Publishers struct {
	Users    []*User  `json:"users"`
	UserKeys []string `json:"userKeys"`
}

// GROUPS //////////////////////////////////////////////////////////////////////////////

type Group struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type GroupCollection struct {
	Results []*Group `json:"results"`
	Start   int      `json:"start"`
	Limit   int      `json:"limit"`
	Size    int      `json:"size"`
}

// RESTRICTIONS ////////////////////////////////////////////////////////////////////////

type Restrictions struct {
	Read   *Restriction `json:"read"`
	Update *Restriction `json:"update"`
}

type Restriction struct {
	Operation string           `json:"operation"`
	Data      *RestrictionData `json:"restrictions"`
}

type RestrictionData struct {
	User  *UserCollection  `json:"user"`
	Group *GroupCollection `json:"group"`
}

// SEARCH //////////////////////////////////////////////////////////////////////////////

type SearchParameters struct {
	CQL                   string   `query:"cql"`
	CQLContext            string   `query:"cqlcontext"`
	Excerpt               string   `query:"excerpt"`
	IncludeArchivedSpaces bool     `query:"includeArchivedSpaces"`
	Expand                []string `query:"expand"`
	Start                 int      `query:"start"`
	Limit                 int      `query:"limit"`
}

type SearchResult struct {
	Results        []*SearchEntity `json:"results"`
	Start          int             `json:"start"`
	Limit          int             `json:"limit"`
	Size           int             `json:"size"`
	TotalSize      int             `json:"totalSize"`
	CQLQuery       string          `json:"cqlQuery"`
	SearchDuration int             `json:"searchDuration"`
}

type SearchEntity struct {
	Title        string `json:"title"`
	Excerpt      string `json:"excerpt"`
	URL          string `json:"url"`
	EntityType   string `json:"entityType"`
	LastModified *Date  `json:"lastModified"`
}

// SPACE ///////////////////////////////////////////////////////////////////////////////

type SpaceParameters struct {
	SpaceKey  []string `query:"spaceKey,unwrap"`
	Type      string   `query:"type"`
	Status    string   `query:"status"`
	Label     string   `query:"label"`
	Favourite bool     `query:"favourite"`
	Depth     string   `query:"depth"`
	Expand    []string `query:"expand"`
	Start     int      `query:"start"`
	Limit     int      `query:"limit"`
}

type Space struct {
	ID   int    `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
	Icon *Icon  `json:"icon"`
	Type string `json:"type"`
}

type Icon struct {
	Path      string `json:"path"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	IsDefault bool   `json:"isDefault"`
}

// USER ////////////////////////////////////////////////////////////////////////////////

type UserParameters struct {
	Key      string   `query:"key"`
	Username string   `query:"username"`
	Expand   []string `query:"expand"`
	Start    int      `query:"start"`
	Limit    int      `query:"limit"`
}

type User struct {
	Type           string `json:"type"`
	Username       string `json:"username"`
	UserKey        string `json:"userKey"`
	ProfilePicture *Icon  `json:"profilePicture"`
	DisplayName    string `json:"displayName"`
}

type UserCollection struct {
	Results []*User `json:"results"`
	Start   int     `json:"start"`
	Limit   int     `json:"limit"`
	Size    int     `json:"size"`
}

// WATCH ///////////////////////////////////////////////////////////////////////////////

type WatchParameters struct {
	Key         string `query:"key"`
	Username    string `query:"username"`
	ContentType string `query:"contentType"`
}

type ListWatchersParameters struct {
	PageID string `query:"pageId"`
}

type WatchStatus struct {
	Watching bool `json:"watching"`
}

type WatchInfo struct {
	PageWatchers  []*Watcher `json:"pageWatchers"`
	SpaceWatchers []*Watcher `json:"spaceWatchers"`
}

type Watcher struct {
	AvatarURL string `json:"avatarUrl"`
	Name      string `json:"name"`
	FullName  string `json:"fullName"`
	Type      string `json:"type"`
	UserKey   string `json:"userKey"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// UnmarshalJSON is custom Date format unmarshaler
func (d *Date) UnmarshalJSON(b []byte) error {
	var err error

	d.Time, err = time.Parse(time.RFC3339, strings.Trim(string(b), "\""))

	return err
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
