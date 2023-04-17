package confluence

// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"regexp"
	"strings"
	"testing"
	"time"

	. "github.com/essentialkaos/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type MyParams struct {
	S  string    `query:"s,respect"`
	I  int       `query:"i,respect"`
	B  bool      `query:"b,respect"`
	BR bool      `query:"br,reverse"`
	BN bool      `query:"bn"`
	DN time.Time `query:"dn"`
}

func (p MyParams) ToQuery() string {
	return paramsToQuery(p)
}

func (p MyParams) Validate() error {
	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

type ConfluenceSuite struct{}

// ////////////////////////////////////////////////////////////////////////////////// //

var _ = Suite(&ConfluenceSuite{})

var tsRegex = regexp.MustCompile(`\&\_\=[0-9]{19}`)

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *ConfluenceSuite) TestParamsEncoding(c *C) {
	var p Parameters

	p = AuditParameters{
		StartDate: time.Date(2018, 1, 1, 0, 0, 0, 0, time.Local),
		EndDate:   time.Date(2018, 2, 15, 12, 30, 0, 0, time.Local),
		Start:     50,
		Limit:     20,
	}

	c.Assert(p.ToQuery(), Equals, "startDate=2018-01-01&endDate=2018-02-15&start=50&limit=20")

	p = CollectionParameters{
		Expand: []string{"test1,test2"},
	}

	c.Assert(p.ToQuery(), Equals, `expand=test1%2Ctest2`)

	p = SpaceParameters{
		SpaceKey:  []string{"TS1", "TS2", "TS3"},
		Favourite: true,
	}

	c.Assert(p.ToQuery(), Equals, "spaceKey=TS1&spaceKey=TS2&spaceKey=TS3&favourite=true")

	p = WatchParameters{}

	c.Assert(p.ToQuery(), Equals, "")

	p = MyParams{BR: true}
	pp := []string{"s=", "i=0", "b=false", "br=false"}

	c.Assert(validateQuery(p.ToQuery(), pp), Equals, true)
}

func (s *ConfluenceSuite) TestTinyLinkGeneration(c *C) {
	api, _ := NewAPI("https://confl.domain.com", AuthBasic{"JohnDoe", "Test1234!"})

	c.Assert(api.GenTinyLink("1477502"), Equals, "https://confl.domain.com/x/fosW")
	c.Assert(api.GenTinyLink("1477627"), Equals, "https://confl.domain.com/x/_4sW")
	c.Assert(api.GenTinyLink("40643836"), Equals, "https://confl.domain.com/x/-CxsAg")
}

func (s *ConfluenceSuite) TestCustomUnmarshalers(c *C) {
	var err error

	d := &Date{}
	err = d.UnmarshalJSON([]byte(`"2013-03-12T10:36:12.602+04:00"`))

	c.Assert(err, IsNil)
	c.Assert(d.Year(), Equals, 2013)
	c.Assert(d.Month(), Equals, time.Month(3))
	c.Assert(d.Day(), Equals, 12)

	t := &Timestamp{}
	err = t.UnmarshalJSON([]byte("1523059214803"))

	c.Assert(err, IsNil)
	c.Assert(t.Year(), Equals, 2018)
	c.Assert(t.Month(), Equals, time.Month(4))
	c.Assert(t.Day(), Equals, 7)

	var e ExtensionPosition
	err = e.UnmarshalJSON([]byte(`"none"`))

	c.Assert(err, IsNil)
	c.Assert(e, Equals, ExtensionPosition(-1))
}

func (s *ConfluenceSuite) TestCalendarIDValidator(c *C) {
	c.Assert(IsValidCalendarID(""), Equals, false)
	c.Assert(IsValidCalendarID("1a72410b-6417-4869-9260-9ec13816e48q"), Equals, false)
	c.Assert(IsValidCalendarID("1a72410b164175486969260f9ec13816e481"), Equals, false)
	c.Assert(IsValidCalendarID("1a72410b-6417-4869-9260-9ec13816e481"), Equals, true)
}

func (s *ConfluenceSuite) TestCalendarParamsEncoding(c *C) {
	p1 := CalendarEventsParameters{
		SubCalendarID:  "1a72410b-6417-4869-9260-9ec13816e481",
		UserTimezoneID: "Etc/UTC",
		Start:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
		End:            time.Date(2020, 1, 2, 12, 30, 45, 0, time.Local),
	}

	pp1 := []string{
		"subCalendarId=1a72410b-6417-4869-9260-9ec13816e481",
		"userTimeZoneId=Etc%2FUTC",
		"start=2020-01-01T00:00:00Z",
		"end=2020-01-02T12:30:45Z",
	}

	q1 := p1.ToQuery()

	c.Assert(validateQuery(q1, pp1), Equals, true)
	c.Assert(tsRegex.MatchString(q1), Equals, true)

	p2 := CalendarsParameters{
		IncludeSubCalendarID: []string{
			"1a72410b-6417-4869-9260-9ec13816e481",
			"1a72410b-6417-4869-9260-9ec13816e482",
		},
		ViewingSpaceKey: "ABC",
		CalendarContext: CALENDAR_CONTEXT_MY,
	}

	pp2 := []string{
		"calendarContext=myCalendars",
		"viewingSpaceKey=ABC",
		"include=1a72410b-6417-4869-9260-9ec13816e481",
		"include=1a72410b-6417-4869-9260-9ec13816e482",
	}

	q2 := p2.ToQuery()

	c.Assert(validateQuery(q2, pp2), Equals, true)
	c.Assert(tsRegex.MatchString(q2), Equals, true)
}

func (s *ConfluenceSuite) TestAuthMethods(c *C) {
	b1 := AuthBasic{"JohnDoe", "Test1234!"}
	b2 := AuthBasic{"", "Test1234!"}
	b3 := AuthBasic{"JohnDoe", ""}

	c.Assert(b1.Encode(), Equals, "Basic Sm9obkRvZTpUZXN0MTIzNCE=")
	c.Assert(b1.Validate(), IsNil)
	c.Assert(b2.Validate(), DeepEquals, ErrEmptyUser)
	c.Assert(b3.Validate(), DeepEquals, ErrEmptyPassword)

	t1 := AuthToken{"TESTVYhExHzKbHzNPCMRmviasXJoUaATysUimxwiWmkr"}
	t2 := AuthToken{""}
	t3 := AuthToken{"TEST"}

	c.Assert(t1.Encode(), Equals, "Bearer TESTVYhExHzKbHzNPCMRmviasXJoUaATysUimxwiWmkr")
	c.Assert(t1.Validate(), IsNil)
	c.Assert(t2.Validate(), DeepEquals, ErrEmptyToken)
	c.Assert(t3.Validate(), DeepEquals, ErrTokenWrongLength)
}

// ////////////////////////////////////////////////////////////////////////////////// //

func validateQuery(query string, parts []string) bool {
	queryParts := strings.Split(query, "&")

LOOP:
	for _, part := range parts {
		for _, qp := range queryParts {
			if part == qp {
				continue LOOP
			}
		}

		return false
	}

	return true
}
