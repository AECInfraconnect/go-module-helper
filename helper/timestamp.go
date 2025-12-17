package helper

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

type Timestamp time.Time

const TimestampLayout = "2006-01-02 15:04:05-07"

/*
------------------------
Timestamp Function
------------------------
*/

func NewTimestampFromString(s string) Timestamp {
	if s == "" {
		return Timestamp(time.Time{})
	}
	loc, _ := time.LoadLocation("Asia/Bangkok")
	ts, err := time.ParseInLocation(TimestampLayout, s, loc)
	if err != nil {
		panic(err)
	}

	return Timestamp(ts)
}

func NewTimestampFromTime(t time.Time) Timestamp {
	loc := time.FixedZone("UTC+7", 7*60*60)
	ts, err := time.Parse(TimestampLayout, t.UTC().Format(TimestampLayout))
	if err != nil {
		panic(err)
	}
	ts = ts.In(loc)

	return Timestamp(ts)
}

func NewTimestampAddDayFromTime(t time.Time, years, months, days int) Timestamp {
	loc := time.FixedZone("UTC+7", 7*60*60)
	ts, err := time.Parse(TimestampLayout, t.UTC().Format(TimestampLayout))
	if err != nil {
		panic(err)
	}
	ts = ts.In(loc).AddDate(years, months, days)

	return Timestamp(ts)
}

func (t Timestamp) Format(f string) string {
	return time.Time(t).Format(f)
}

func (t Timestamp) String() string {
	return t.Format(TimestampLayout)
}

func (t Timestamp) ToTime() time.Time {
	return time.Time(t)
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	ts, err := time.Parse(TimestampLayout, s)
	if err != nil {
		return err
	}
	*t = Timestamp(ts)

	return nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Format(TimestampLayout))
}

func (t Timestamp) Value() (driver.Value, error) {
	if t == (Timestamp{}) {
		return nil, nil
	}

	return t.String(), nil
}

func (t *Timestamp) Scan(v interface{}) error {
	if v == nil {
		*t = Timestamp(time.Time{})
		return nil
	}
	ts, ok := v.(time.Time)
	if !ok {
		return nil
	}

	*t = Timestamp(ts)
	return nil
}
