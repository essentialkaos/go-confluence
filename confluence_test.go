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

	c.Assert(p.ToQuery(), Equals, "spaceKey=TS1&spaceKey=TS2&spaceKey=TS3&favourite=1")
}

func (s *ConfluenceSuite) TestCustomUnmarshalers(c *C) {
	d := &Date{}
	d.UnmarshalJSON([]byte("\"2013-03-12T10:36:12.602+04:00\""))

	c.Assert(d.Year(), Equals, 2013)
	c.Assert(d.Month(), Equals, time.Month(3))
	c.Assert(d.Day(), Equals, 12)

	t := &Timestamp{}
	t.UnmarshalJSON([]byte("1523059214803"))

	c.Assert(t.Year(), Equals, 2018)
	c.Assert(t.Month(), Equals, time.Month(4))
	c.Assert(t.Day(), Equals, 7)
}
