package confluence

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2022 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"errors"
	"regexp"
	"time"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const _REST_BASE = "/rest/calendar-services/1.0"

// ////////////////////////////////////////////////////////////////////////////////// //

// Calendar context
const (
	CALENDAR_CONTEXT_MY    = "myCalendars"
	CALENDAR_CONTEXT_SPACE = "spaceCalendars"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// CalendarEventsParameters contains request params for events from Team Calendars API
type CalendarEventsParameters struct {
	SubCalendarID  string    `query:"subCalendarId"`
	UserTimezoneID string    `query:"userTimeZoneId"`
	Start          time.Time `query:"start,timedate"`
	End            time.Time `query:"end,timedate"`

	timestamp int64 `query:"_"`
}

// CalendarsParameters contains request params for calendars from Team Calendars API
type CalendarsParameters struct {
	IncludeSubCalendarID []string `query:"include,unwrap"`
	CalendarContext      string   `query:"calendarContext"`
	ViewingSpaceKey      string   `query:"viewingSpaceKey"`

	timestamp int64 `query:"_"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// CalendarEventCollection contains slice with events
type CalendarEventCollection struct {
	Events  []*CalendarEvent `json:"events"`
	Success bool             `json:"success"`
}

// CalendarCollection contains slice with calendars
type CalendarCollection struct {
	Calendars []*Calendar `json:"payload"`
	Success   bool        `json:"success"`
}

// Calendar represents Team Calendars calendar
type Calendar struct {
	UsersPermittedToView      []*PermsUser `json:"usersPermittedToView"`
	UsersPermittedToEdit      []*PermsUser `json:"usersPermittedToEdit"`
	GroupsPermittedToView     []string     `json:"groupsPermittedToView"`
	GroupsPermittedToEdit     []string     `json:"groupsPermittedToEdit"`
	Warnings                  []string     `json:"warnings"`
	ChildSubCalendars         []*Calendar  `json:"childSubCalendars"`
	SubscriberCount           int          `json:"subscriberCount"`
	SubCalendar               *SubCalendar `json:"subCalendar"`
	ReminderMe                bool         `json:"reminderMe"`
	IsHidden                  bool         `json:"hidden"`
	IsEditable                bool         `json:"editable"`
	IsReloadable              bool         `json:"reloadable"`
	IsDeletable               bool         `json:"deletable"`
	IsEventsHidden            bool         `json:"eventsHidden"`
	IsWatchedViaContent       bool         `json:"watchedViaContent"`
	IsAdministrable           bool         `json:"administrable"`
	IsWatched                 bool         `json:"watched"`
	IsEventsViewable          bool         `json:"eventsViewable"`
	IsEventsEditable          bool         `json:"eventsEditable"`
	IsSubscribedByCurrentUser bool         `json:"subscribedByCurrentUser"`
}

// SubCalendar represents Team Calendars sub-calendar
type SubCalendar struct {
	DisableEventTypes        []string             `json:"disableEventTypes"`
	CustomEventTypes         []*CustomEventType   `json:"customEventTypes"`
	SanboxEventTypeReminders []*EventTypeReminder `json:"sanboxEventTypeReminders"`
	Creator                  string               `json:"creator"`
	TypeKey                  string               `json:"typeKey"`
	Color                    string               `json:"color"`
	TimeZoneID               string               `json:"timeZoneId"`
	Description              string               `json:"description"`
	Type                     string               `json:"type"`
	SpaceKey                 string               `json:"spaceKey"`
	SpaceName                string               `json:"spaceName"`
	Name                     string               `json:"name"`
	ID                       string               `json:"id"`
	IsWatchable              bool                 `json:"watchable"`
	IsEventInviteesSupported bool                 `json:"eventInviteesSupported"`
	IsRestrictable           bool                 `json:"restrictable"`
}

// CustomEventType contains info about custom event type
type CustomEventType struct {
	Created             string `json:"created"`
	Icon                string `json:"icon"`
	PeriodInMins        int    `json:"periodInMins"`
	CustomEventTypeID   string `json:"customEventTypeId"`
	Title               string `json:"title"`
	ParentSubCalendarID string `json:"parentSubCalendarId"`
}

// EventTypeReminder contains info about event reminder
type EventTypeReminder struct {
	EventTypeID       string `json:"eventTypeId"`
	PeriodInMins      int    `json:"periodInMins"`
	IsCustomEventType bool   `json:"isCustomEventType"`
}

// CalendarEvent represents Team Calendars event
type CalendarEvent struct {
	Invitees              []*CalendarUser `json:"invitees"`
	WorkingURL            string          `json:"workingUrl"`
	Description           string          `json:"description"`
	ClassName             string          `json:"className"`
	ShortTitle            string          `json:"shortTitle"`
	Title                 string          `json:"title"`
	EventType             string          `json:"eventType"`
	ID                    string          `json:"id"`
	CustomEventTypeID     string          `json:"customEventTypeId"`
	SubCalendarID         string          `json:"subCalendarId"`
	IconURL               string          `json:"iconUrl"`
	IconLink              string          `json:"iconLink"`
	MediumIconURL         string          `json:"mediumIconUrl"`
	BackgroundColor       string          `json:"backgroundColor"`
	BorderColor           string          `json:"borderColor"`
	TextColor             string          `json:"textColor"`
	ColorScheme           string          `json:"colorScheme"`
	Where                 string          `json:"where"`
	FormattedStartDate    string          `json:"confluenceFormattedStartDate"`
	Start                 *Date           `json:"start"`
	End                   *Date           `json:"end"`
	OriginalStartDateTime *Date           `json:"originalStartDateTime"`
	OriginalEndDateTime   *Date           `json:"originalEndDateTime"`
	IsExpandDates         bool            `json:"expandDates"`
	IsEditable            bool            `json:"editable"`
	IsAllDay              bool            `json:"allDay"`
}

// CalendarUser represents Team Calendars user
type CalendarUser struct {
	DisplayName   string `json:"displayName"`
	Name          string `json:"name"`
	ID            string `json:"id"`
	Type          string `json:"type"`
	AvatarIconURL string `json:"avatarIconUrl"`
	Email         string `json:"email"`
}

// PermsUser represents Team Calendars permissions user
type PermsUser struct {
	AvatarURL   string `json:"avatarUrl"`
	Name        string `json:"name"`
	DisplayName string `json:"fullName"`
	Key         string `json:"id"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

var idValidationRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// ////////////////////////////////////////////////////////////////////////////////// //

// GetCalendarEvents fetch events from given calendar
func (api *API) GetCalendarEvents(params CalendarEventsParameters) (*CalendarEventCollection, error) {
	result := &CalendarEventCollection{}
	statusCode, err := api.doRequest(
		"GET", _REST_BASE+"/calendar/events.json",
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	if statusCode == 403 {
		return nil, ErrNoPerms
	}

	return result, nil
}

func (api *API) GetCalendars(params CalendarsParameters) (*CalendarCollection, error) {
	result := &CalendarCollection{}
	statusCode, err := api.doRequest(
		"GET", _REST_BASE+"/calendar/subcalendars.json",
		params, result, nil,
	)

	if err != nil {
		return nil, err
	}

	if statusCode == 403 {
		return nil, ErrNoPerms
	}

	return result, nil
}

// IsValidCalendarID validates calendar ID
func IsValidCalendarID(id string) bool {
	return idValidationRegex.MatchString(id)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Validate validates parameters
func (p CalendarEventsParameters) Validate() error {
	switch {
	case p.SubCalendarID == "":
		return errors.New("SubCalendarID is mandatory and must be set")

	case !IsValidCalendarID(p.SubCalendarID):
		return errors.New("SubCalendarID contains invalid calendar ID")

	case p.UserTimezoneID == "":
		return errors.New("UserTimezoneID is mandatory and must be set")

	case p.Start.IsZero():
		return errors.New("Start is mandatory and must be set")

	case p.End.IsZero():
		return errors.New("End is mandatory and must be set")
	}

	return nil
}

// Validate validates parameters
func (p CalendarsParameters) Validate() error {
	if p.CalendarContext == "" {
		return errors.New("CalendarContext is mandatory and must be set")
	}

	if p.CalendarContext == CALENDAR_CONTEXT_MY {
		return nil
	}

	switch {
	case len(p.IncludeSubCalendarID) == 0:
		return errors.New("IncludeSubCalendarID is mandatory and must be set")

	case p.ViewingSpaceKey == "":
		return errors.New("ViewingSpaceKey is mandatory and must be set")
	}

	for _, id := range p.IncludeSubCalendarID {
		if id == "" {
			return errors.New("IncludeSubCalendarID is mandatory and must be set")
		}

		if !IsValidCalendarID(id) {
			return errors.New("IncludeSubCalendarID contains invalid calendar ID")
		}
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// ToQuery convert params to URL query
func (p CalendarEventsParameters) ToQuery() string {
	p.timestamp = time.Now().UnixNano()
	return paramsToQuery(p)
}

// ToQuery convert params to URL query
func (p CalendarsParameters) ToQuery() string {
	p.timestamp = time.Now().UnixNano()
	return paramsToQuery(p)
}
