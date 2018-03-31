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

const (
	NAME    = "Go-Confluence"
	VERSION = "1.0.0"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type API struct {
	// Client is client for http requests
	Client *fasthttp.Client

	url                string // Confluence URL
	basicAuth          string // Basic auth
	clientInitComplete bool   // client init flag
}

// ////////////////////////////////////////////////////////////////////////////////// //

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

var clientInitComplete bool

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
func (api *API) GetSpaces(params SpaceParameters) ([]*Space, error) {
	result := &struct {
		Results []*Space `json:"results"`
		Size    int      `json:"size"`
	}{}

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

	return result.Results, nil
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

// GetContent fetch the content in this given space
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#space-contents
func (api *API) GetContent(spaceKey string, params SpaceParameters) (*Contents, error) {
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

// GetContentWithType fetch the content in this given space with the given type
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#space-contentsWithType
func (api *API) GetContentWithType(spaceKey, contentType string, params SpaceParameters) (*Contents, error) {
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

// GetAnonymous fetch information about the how anonymous is represented in confluence
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#user-getAnonymous
func (api *API) GetAnonymous() (*User, error) {
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

// GetCurrent fetch information about the current logged in user
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#user-getCurrent
func (api *API) GetCurrent(params Parameters) (*User, error) {
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

// GetGroups fetch collection of groups that the given user is a member of
// https://docs.atlassian.com/ConfluenceServer/rest/6.8.0/#user-getGroups
func (api *API) GetGroups(params UserParameters) ([]*Group, error) {
	result := &struct {
		Results []*Group `json:"results"`
		Size    int      `json:"size"`
	}{}

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

	return result.Results, nil
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

// acquireRequest acquire new request with given params
func (api *API) acquireRequest(method, uri string, params Parameters) *fasthttp.Request {
	req := fasthttp.AcquireRequest()
	query := params.ToQuery()

	req.SetRequestURI(api.url + uri)

	// Set query if params can be encoded as query
	if query != "" {
		req.URI().SetQueryString(query)
	}

	// TODO: DEBUG / REMOVE ON RELEASE
	fmt.Println("→", uri, "»", query)

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
