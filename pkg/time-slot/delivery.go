package time_slot

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type (
	// Hour ...
	Hour uint

	// Weekday ...
	Weekday = time.Weekday

	//Delivery ...
	Delivery struct {
		WDay Weekday
		From Hour
		To   Hour
	}
)

// ConstructDelivery ...
func ConstructDelivery(wd Weekday, from, to uint) Delivery {
	return Delivery{
		WDay: wd,
		From: Hour(from),
		To:   Hour(to),
	}
}

func (ts Delivery) String() string {
	return fmt.Sprintf("%s %s - %s", ts.WDay, ts.From, ts.To)
}

// MarshalJSON ...
func (ts Delivery) MarshalJSON() ([]byte, error) {
	const api = "Delivery.UnmarshalJSON"

	var b bytes.Buffer
	e := b.WriteByte(byte('"'))
	if e == nil {
		_, e = b.WriteString(ts.String())
	}
	if e == nil {
		e = b.WriteByte(byte('"'))
	}
	return b.Bytes(), errors.Wrap(e, api)
}

// UnmarshalJSON ...
func (ts *Delivery) UnmarshalJSON(data []byte) error {
	const api = "Delivery.UnmarshalJSON"

	allIndexes := deliveryRE.FindAllSubmatchIndex(data, -1)
	if len(allIndexes) == 0 {
		return errors.Errorf("%s: wrong incoming data %q", api, string(data))
	}

	sub := allIndexes[0]
	if len(sub) < 8 {
		return errors.Errorf("%s: wrong incoming data %q", api, string(data))
	}

	var ok bool
	s := data[sub[2]:sub[3]]
	if ts.WDay, ok = s2wd[string(s)]; !ok {
		return errors.Errorf("%s: wrong incoming data %q", api, string(data))
	}
	var e error
	s = data[sub[4]:sub[5]]
	if e = ts.From.FromString(s); e == nil {
		s = data[sub[6]:sub[7]]
		e = ts.To.FromString(s)
	}
	return errors.Wrapf(e, "%s: wrong incoming data %q", api, string(data))
}

// IsIn ...
func (ts Delivery) IsIn(ts1 Delivery) bool {
	return ts.WDay == ts1.WDay &&
		ts.From >= ts1.From &&
		ts.To <= ts1.To
}

// FromString ...
func (h *Hour) FromString(data []byte) error {
	const (
		h24 = iota
		am
		pm
	)
	var amPm = h24
	var val uint64
	var e error
	data = bytes.ToUpper(data)
	if bytes.HasSuffix(data, []byte("AM")) {
		amPm = am
		data = data[:len(data)-2]
	} else if bytes.HasSuffix(data, []byte("PM")) {
		amPm = pm
		data = data[:len(data)-2]
	}
	if val, e = strconv.ParseUint(string(data), 10, 0); e != nil {
		return e
	}
	switch amPm {
	case am, pm:
		if val > 12 {
			return errors.New("bad 12h value")
		}
		if amPm == am && val == 12 {
			val = 0
		} else if amPm == pm && val < 12 {
			val += 12
		}
	default:
		if val >= 24 {
			return errors.New("bad 24h value")
		}
	}
	*h = Hour(val)
	return nil
}

func (h Hour) String() string {
	if h == 0 {
		return "12AM"
	}
	if h == 12 {
		return "12PM"
	}
	if h < 12 {
		return fmt.Sprintf("%dAM", uint(h))
	}
	return fmt.Sprintf("%dPM", uint(h)-12)
}

//MarshalJSON ...
func (h Hour) MarshalJSON() ([]byte, error) {
	const api = "Hour.UnmarshalJSON"

	var b bytes.Buffer
	e := b.WriteByte(byte('"'))
	if e == nil {
		_, e = b.WriteString(h.String())
	}
	if e == nil {
		e = b.WriteByte(byte('"'))
	}
	return b.Bytes(), errors.Wrap(e, api)
}

var (
	s2wd = map[string]time.Weekday{
		time.Sunday.String():    time.Sunday,
		time.Monday.String():    time.Monday,
		time.Tuesday.String():   time.Tuesday,
		time.Wednesday.String(): time.Wednesday,
		time.Thursday.String():  time.Thursday,
		time.Friday.String():    time.Friday,
		time.Saturday.String():  time.Saturday,
	}

	deliveryRE = regexp.MustCompile(`(?i)"(\w+)\s*(\d*(?:AM|PM))\s*-\s*(\d+(?:AM|PM))"`)
)
