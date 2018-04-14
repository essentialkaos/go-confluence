package confluence

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2018 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/erikdubbelboer/fasthttp"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// API is Confluence API struct
type API struct {
	Client *fasthttp.Client // Client is client for http requests

	url       string // confluence URL
	basicAuth string // basic auth
}

// ////////////////////////////////////////////////////////////////////////////////// //

// API errors
var (
	ErrInitEmptyURL      = errors.New("URL can't be empty")
	ErrInitEmptyUser     = errors.New("User can't be empty")
	ErrInitEmptyPassword = errors.New("Password can't be empty")
	ErrNoPerms           = errors.New("User does not have permission to use confluence")
	ErrQueryError        = errors.New("Query cannot be parsed")
	ErrNoContent         = errors.New("There is no content with the given id, or if the calling user does not have permission to view the content")
	ErrNoSpace           = errors.New("There is no space with the given key, or if the calling user does not have permission to view the space")
	ErrNoUserPerms       = errors.New("User does not have permission to view users")
	ErrNoUserFound       = errors.New("User with the given username or userkey does not exist")
)

// ////////////////////////////////////////////////////////////////////////////////// //

// NewAPI create new API struct
func NewAPI(url, username, password string) (*API, error) {
	switch {
	case url == "":
		return nil, ErrInitEmptyURL
	case username == "":
		return nil, ErrInitEmptyUser
	case password == "":
		return nil, ErrInitEmptyPassword
	}

	return &API{
		Client: &fasthttp.Client{
			Name:                getUserAgent("", ""),
			MaxIdleConnDuration: 5 * time.Second,
			ReadTimeout:         3 * time.Second,
			WriteTimeout:        3 * time.Second,
			MaxConnsPerHost:     150,
		},

		url:       url,
		basicAuth: genBasicAuthHeader(username, password),
	}, nil
}

// SetUserAgent set user-agent string based on app name and version
func (api *API) SetUserAgent(app, version string) {
	api.Client.Name = getUserAgent(app, version)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// GetAuditRecords fetch a list of AuditRecord instances dating back to a certain time
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#audit-getAuditRecords
func (api *API) GetAuditRecords(params AuditParameters) (*AuditRecordCollection, error) {
	result := &AuditRecordCollection{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/audit",
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	}

	return result, nil
}

// GetAuditRecordsSince fetch a list of AuditRecord instances dating back to a certain time
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#audit-getAuditRecords
func (api *API) GetAuditRecordsSince(params AuditSinceParameters) (*AuditRecordCollection, error) {
	result := &AuditRecordCollection{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/audit/since",
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	}

	return result, nil
}

// GetAuditRetention fetch the current retention period
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#audit-getRetentionPeriod
func (api *API) GetAuditRetention() (*AuditRetentionInfo, error) {
	result := &AuditRetentionInfo{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/audit/retention",
		EmptyParameters{}, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	}

	return result, nil
}

// GetContent fetch list of Content
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#content-getContent
func (api *API) GetContent(params ContentParameters) (*ContentColletion, error) {
	result := &ContentColletion{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/content",
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	case 404:
		return nil, ErrNoContent
	}

	return result, nil
}

// GetContentByID fetch a piece of Content
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#content-getContentById
func (api *API) GetContentByID(contentID string, params ContentIDParameters) (*Content, error) {
	result := &Content{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/content/"+contentID,
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	case 404:
		return nil, ErrNoContent
	}

	return result, nil
}

// GetContentHistory fetch the history of a particular piece of content
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#content-getHistory
func (api *API) GetContentHistory(contentID string, params ExpandParameters) (*History, error) {
	result := &History{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/content/"+contentID+"/history",
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	case 404:
		return nil, ErrNoContent
	}

	return result, nil
}

// GetContentChildren fetch a map of the direct children of a piece of Content
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#content/{id}/child-children
func (api *API) GetContentChildren(contentID string, params ChildrenParameters) (*Contents, error) {
	result := &Contents{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/content/"+contentID+"/child",
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	case 404:
		return nil, ErrNoContent
	}

	return result, nil
}

// GetContentChildrenByType the direct children of a piece of Content, limited to a single child type
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#content/{id}/child-childrenOfType
func (api *API) GetContentChildrenByType(contentID, contentType string, params ChildrenParameters) (*ContentColletion, error) {
	result := &ContentColletion{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/content/"+contentID+"/child/"+contentType,
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	case 404:
		return nil, ErrNoContent
	}

	return result, nil
}

// GetContentComments fetch the comments of a content
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#content/{id}/child-commentsOfContent
func (api *API) GetContentComments(contentID string, params ChildrenParameters) (*ContentColletion, error) {
	result := &ContentColletion{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/content/"+contentID+"/child/comment",
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	case 404:
		return nil, ErrNoContent
	}

	return result, nil
}

// GetAttachments fetch list of attachment Content entities within a single container
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#content/{id}/child/attachment-getAttachments
func (api *API) GetAttachments(contentID string, params AttachmentParameters) (*ContentColletion, error) {
	result := &ContentColletion{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/content/"+contentID+"/child/attachment",
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	case 404:
		return nil, ErrNoContent
	}

	return result, nil
}

// GetDescendants fetch a map of the descendants of a piece of Content
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#content/{id}/descendant-descendants
func (api *API) GetDescendants(contentID string, params ExpandParameters) (*Contents, error) {
	result := &Contents{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/content/"+contentID+"/descendant",
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	case 404:
		return nil, ErrNoContent
	}

	return result, nil
}

// GetDescendantsOfType fetch the direct descendants of a piece of Content, limited to a single descendant type
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#content/{id}/descendant-descendantsOfType
func (api *API) GetDescendantsOfType(contentID, descType string, params ExpandParameters) (*ContentColletion, error) {
	result := &ContentColletion{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/content/"+contentID+"/descendant/"+descType,
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	case 404:
		return nil, ErrNoContent
	}

	return result, nil
}

// GetLabels fetch the list of labels on a piece of Content
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#content/{id}/label-labels
func (api *API) GetLabels(contentID string, params LabelParameters) (*LabelCollection, error) {
	result := &LabelCollection{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/content/"+contentID+"/label",
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	case 404:
		return nil, ErrNoContent
	}

	return result, nil
}

// GetRestrictionsByOperation fetch info about all restrictions by operation
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#content/{id}/restriction-byOperation
func (api *API) GetRestrictionsByOperation(contentID string, params ExpandParameters) (*Restrictions, error) {
	result := &Restrictions{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/content/"+contentID+"/restriction/byOperation",
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	}

	return result, nil
}

// GetRestrictionsForOperation fetch info about all restrictions of given operation
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#content/{id}/restriction-forOperation
func (api *API) GetRestrictionsForOperation(contentID, operation string, params CollectionParameters) (*Restriction, error) {
	result := &Restriction{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/content/"+contentID+"/restriction/byOperation/"+operation,
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	}

	return result, nil
}

// GetGroups fetch collection of user groups
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#group-getGroups
func (api *API) GetGroups(params CollectionParameters) (*GroupCollection, error) {
	result := &GroupCollection{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/group",
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	}

	return result, nil
}

// GetGroup fetch the user group with the group name
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#group-getGroup
func (api *API) GetGroup(groupName string, params ExpandParameters) (*Group, error) {
	result := &Group{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/group/"+groupName,
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	}

	return result, nil
}

// GetGroupMembers fetch a collection of users in the given group
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#group-getMembers
func (api *API) GetGroupMembers(groupName string, params CollectionParameters) (*UserCollection, error) {
	result := &UserCollection{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/group/"+groupName+"/member",
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	}

	return result, nil
}

// Search search for entities in Confluence using the Confluence Query Language (CQL)
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#search-search
func (api *API) Search(params SearchParameters) (*SearchResult, error) {
	result := &SearchResult{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/search",
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 400:
		return nil, ErrQueryError
	case 403:
		return nil, ErrNoPerms
	}

	return result, nil
}

// GetSpaces fetch information about a number of spaces
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#space-spaces
func (api *API) GetSpaces(params SpaceParameters) (*SpaceCollection, error) {
	result := &SpaceCollection{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/space",
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	}

	return result, nil
}

// GetSpace fetch information about a space
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#space-space
func (api *API) GetSpace(spaceKey string, params Parameters) (*Space, error) {
	result := &Space{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/space/"+spaceKey,
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	case 404:
		return nil, ErrNoSpace
	}

	return result, nil
}

// GetSpaceContent fetch the content in this given space
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#space-contents
func (api *API) GetSpaceContent(spaceKey string, params SpaceParameters) (*Contents, error) {
	result := &Contents{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/space/"+spaceKey+"/content",
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	case 404:
		return nil, ErrNoContent
	}

	return result, nil
}

// GetSpaceContentWithType fetch the content in this given space with the given type
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#space-contentsWithType
func (api *API) GetSpaceContentWithType(spaceKey, contentType string, params SpaceParameters) (*Contents, error) {
	result := &Contents{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/space/"+spaceKey+"/content/"+contentType,
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	case 404:
		return nil, ErrNoContent
	}

	return result, nil
}

// GetUser fetch information about a user identified by either user key or username
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#user-getUser
func (api *API) GetUser(params UserParameters) (*User, error) {
	result := &User{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/user",
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoUserPerms
	case 404:
		return nil, ErrNoUserFound
	}

	return result, nil
}

// GetAnonymousUser fetch information about the how anonymous is represented in confluence
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#user-getAnonymous
func (api *API) GetAnonymousUser() (*User, error) {
	result := &User{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/user/anonymous",
		EmptyParameters{}, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	}

	return result, nil
}

// GetCurrentUser fetch information about the current logged in user
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#user-getCurrent
func (api *API) GetCurrentUser(params ExpandParameters) (*User, error) {
	result := &User{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/user/current",
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	}

	return result, nil
}

// GetUserGroups fetch collection of groups that the given user is a member of
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#user-getGroups
func (api *API) GetUserGroups(params UserParameters) (*GroupCollection, error) {
	result := &GroupCollection{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/user/memberof",
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	}

	return result, nil
}

// IsWatchingContent fetch information about whether a user is watching a specified content
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#user/watch-isWatchingContent
func (api *API) IsWatchingContent(contentID string, params WatchParameters) (*WatchStatus, error) {
	result := &WatchStatus{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/user/watch/content/"+contentID,
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	case 404:
		return nil, ErrNoContent
	}

	return result, nil
}

// IsWatchingSpace fetch information about whether a user is watching a specified space
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#user/watch-isWatchingSpace
func (api *API) IsWatchingSpace(spaceKey string, params WatchParameters) (*WatchStatus, error) {
	result := &WatchStatus{}
	statusCode, err := api.doRequest(
		"GET", "/rest/api/user/watch/space/"+spaceKey,
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	case 404:
		return nil, ErrNoSpace
	}

	return result, nil
}

// ListWatchers fetch information about all watcher of given page
func (api *API) ListWatchers(params ListWatchersParameters) (*WatchInfo, error) {
	result := &WatchInfo{}
	statusCode, err := api.doRequest(
		"GET", "/json/listwatchers.action",
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 403:
		return nil, ErrNoPerms
	case 404:
		return nil, ErrNoSpace
	}

	return result, nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// codebeat:disable[ARITY]

// doRequest create and execute request
func (api *API) doRequest(method, uri string, params Parameters, result, body interface{}) (int, error) {
	req := api.acquireRequest(method, uri, params)
	resp := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	if body != nil {
		bodyData, err := json.Marshal(body)

		if err != nil {
			return -1, err
		}

		req.SetBody(bodyData)
	}

	err := api.Client.Do(req, resp)

	if err != nil {
		return -1, err
	}

	statusCode := resp.StatusCode()

	if statusCode != 200 || result == nil {
		return statusCode, nil
	}

	err = json.Unmarshal(resp.Body(), result)

	return statusCode, err
}

// codebeat:enable[ARITY]

// acquireRequest acquire new request with given params
func (api *API) acquireRequest(method, uri string, params Parameters) *fasthttp.Request {
	req := fasthttp.AcquireRequest()
	query := params.ToQuery()

	req.SetRequestURI(api.url + uri)

	// Set query if params can be encoded as query
	if query != "" {
		req.URI().SetQueryString(query)
	}

	if method != "GET" {
		req.Header.SetMethod(method)
	}

	// Set auth header
	req.Header.Set("Authorization", "Basic "+api.basicAuth)

	return req
}

// ////////////////////////////////////////////////////////////////////////////////// //

// getUserAgent generate user-agent string for client
func getUserAgent(app, version string) string {
	if app != "" && version != "" {
		return fmt.Sprintf(
			"%s/%s %s/%s (go; %s; %s-%s)",
			app, version, NAME, VERSION, runtime.Version(),
			runtime.GOARCH, runtime.GOOS,
		)
	}

	return fmt.Sprintf(
		"%s/%s (go; %s; %s-%s)",
		NAME, VERSION, runtime.Version(),
		runtime.GOARCH, runtime.GOOS,
	)
}

// genBasicAuthHeader generate basic auth header
func genBasicAuthHeader(username, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

}
