package confluence

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2018 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"errors"
	"strconv"
	"time"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const _CALENDAR_TIME_FORMAT = "2006-01-02T15:04:05Z"

// ////////////////////////////////////////////////////////////////////////////////// //

type CalendarParameters struct {
	SubCalendarID  string
	UserTimezoneID string
	Start          time.Time
	End            time.Time
}

type CalendarEventCollection struct {
	Success bool             `json:"success"`
	Events  []*CalendarEvent `json:"events"`
}

type CalendarEvent struct {
	WorkingURL        string `json:"workingUrl"`
	Description       string `json:"description"`
	ClassName         string `json:"className"`
	ShortTitle        string `json:"shortTitle"`
	Title             string `json:"title"`
	EventType         string `json:"eventType"`
	ID                string `json:"id"`
	CustomEventTypeID string `json:"customEventTypeId"`
	SubCalendarID     string `json:"subCalendarId"`
	ExpandDates       bool   `json:"expandDates"`
	Editable          bool   `json:"editable"`
	AllDay            bool   `json:"allDay"`

	Start                 *Date `json:"start"`
	End                   *Date `json:"end"`
	OriginalStartDateTime *Date `json:"originalStartDateTime"`
	OriginalEndDateTime   *Date `json:"originalEndDateTime"`

	IconURL       string `json:"iconUrl"`
	MediumIconURL string `json:"mediumIconUrl"`

	BackgroundColor string `json:"backgroundColor"`
	BorderColor     string `json:"borderColor"`
	TextColor       string `json:"textColor"`
	ColorScheme     string `json:"colorScheme"`

	Invitees []*CalendarUser `json:"invitees"`
}

type CalendarUser struct {
	DisplayName   string `json:"displayName"`
	Name          string `json:"name"`
	ID            string `json:"id"`
	Type          string `json:"type"`
	AvatarIconURL string `json:"avatarIconUrl"`
	Email         string `json:"email"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

var ErrNoCalendarID = errors.New("Sub calendar ID must be defined")

// ////////////////////////////////////////////////////////////////////////////////// //

// GetCalendarEvents fetch events from given calendar
func (api *API) GetCalendarEvents(params CalendarParameters) (*CalendarEventCollection, error) {
	if len(params.SubCalendarID) != 36 {
		return nil, ErrNoCalendarID
	}

	result := &CalendarEventCollection{}
	statusCode, err := api.doRequest(
		"GET", "/rest/calendar-services/1.0/calendar/events.json",
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

// ////////////////////////////////////////////////////////////////////////////////// //

// codebeat:disable[CYCLO]

// ToQuery convert params to URL query
func (p CalendarParameters) ToQuery() string {
	result := "subCalendarId=" + p.SubCalendarID + "&"

	if p.UserTimezoneID == "" {
		result += "userTimeZoneId=Etc/UTC&"
	} else {
		result += "userTimeZoneId=" + p.UserTimezoneID + "&"
	}

	now := time.Now()

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

// codebeat:enable[CYCLO]
