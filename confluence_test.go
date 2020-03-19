package confluence

// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"testing"
	"time"

	. "pkg.re/check.v1"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

type ConfluenceSuite struct{}

// ////////////////////////////////////////////////////////////////////////////////// //

var _ = Suite(&ConfluenceSuite{})

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
}

func (s *ConfluenceSuite) TestTinyLinkGeneration(c *C) {
	api, _ := NewAPI("https://confl.domain.com", "user", "pass")

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
