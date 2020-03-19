package confluence

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2020 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const _CALENDAR_TIME_FORMAT = "2006-01-02T15:04:05Z"

const _REST_BASE = "/rest/calendar-services/1.0"

// ////////////////////////////////////////////////////////////////////////////////// //

// Calendar context
const (
	CALENDAR_CONTEXT_MY    = "myCalendars"
	CALENDAR_CONTEXT_SPACE = "spaceCalendars"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// EventsParameters contains request params for events from Team Calendars API
type EventsParameters struct {
	SubCalendarID  string
	UserTimezoneID string
	Start          time.Time
	End            time.Time
}

// CalendarsParameters contains request params for calendars from Team Calendars API
type CalendarsParameters struct {
	IncludeSubCalendarID []string
	CalendarContext      string
	ViewingSpaceKey      string
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
	UsersPermittedToView      []string     `json:"usersPermittedToView"`
	UsersPermittedToEdit      []string     `json:"usersPermittedToEdit"`
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
	MediumIconURL         string          `json:"mediumIconUrl"`
	BackgroundColor       string          `json:"backgroundColor"`
	BorderColor           string          `json:"borderColor"`
	TextColor             string          `json:"textColor"`
	ColorScheme           string          `json:"colorScheme"`
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

// ////////////////////////////////////////////////////////////////////////////////// //

// ErrWrongIDFormat returns if sub-calendar ID has the wrong format
var ErrWrongIDFormat = errors.New("Sub-calendar ID has the wrong format")

// ErrNoID returns if sub-calendar ID is not defined
var ErrNoID = errors.New("Sub-calendar ID is mandatory")

// ////////////////////////////////////////////////////////////////////////////////// //

var idValidationRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// ////////////////////////////////////////////////////////////////////////////////// //

// GetCalendarEvents fetch events from given calendar
func (api *API) GetCalendarEvents(params EventsParameters) (*CalendarEventCollection, error) {
	if params.SubCalendarID == "" {
		return nil, ErrNoID
	}

	if !IsValidCalendarID(params.SubCalendarID) {
		return nil, ErrWrongIDFormat
	}

	result := &CalendarEventCollection{}
	statusCode, err := api.doRequest(
		"GET", _REST_BASE+"/calendar/events.json",
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

func (api *API) GetCalendars(params CalendarsParameters) (*CalendarCollection, error) {
	for _, id := range params.IncludeSubCalendarID {
		if id == "" {
			return nil, ErrNoID
		}

		if !IsValidCalendarID(id) {
			return nil, ErrWrongIDFormat
		}
	}

	result := &CalendarCollection{}
	statusCode, err := api.doRequest(
		"GET", _REST_BASE+"/calendar/subcalendars.json",
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

// IsValidCalendarID validates calendar ID
func IsValidCalendarID(id string) bool {
	return idValidationRegex.MatchString(id)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// codebeat:disable[CYCLO]

// ToQuery convert params to URL query
func (p EventsParameters) ToQuery() string {
	now := time.Now()
	result := "subCalendarId=" + p.SubCalendarID + "&"

	if p.UserTimezoneID == "" {
		result += "userTimeZoneId=Etc/UTC&"
	} else {
		result += "userTimeZoneId=" + p.UserTimezoneID + "&"
	}

	if p.Start.IsZero() {
		result += "start=" + now.Add(time.Hour*-720).Format(_CALENDAR_TIME_FORMAT) + "&"
	} else {
		result += "start=" + p.Start.Format(_CALENDAR_TIME_FORMAT) + "&"
	}

	if p.End.IsZero() {
		result += "end=" + now.Add(time.Hour*720).Format(_CALENDAR_TIME_FORMAT) + "&"
	} else {
		result += "end=" + p.End.Format(_CALENDAR_TIME_FORMAT) + "&"
	}

	return result + "_=" + strconv.FormatInt(now.UnixNano(), 10)
}

// ToQuery convert params to URL query
func (p CalendarsParameters) ToQuery() string {
	var result string

	now := time.Now()

	if p.CalendarContext != "" {
		result = "calendarContext=" + p.CalendarContext + "&"
	} else {
		result = "calendarContext=" + CALENDAR_CONTEXT_MY + "&"
	}

	if p.ViewingSpaceKey != "" {
		result += "viewingSpaceKey=" + p.ViewingSpaceKey + "&"
	}

	for _, id := range p.IncludeSubCalendarID {
		result += "include=" + id + "&"
	}

	return result + "_=" + strconv.FormatInt(now.UnixNano(), 10)
}

// codebeat:enable[CYCLO]
