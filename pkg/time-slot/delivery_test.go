package time_slot

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTsMarshalUnmarshal(t *testing.T) {
	type ts struct {
		TS Delivery `json:"ts"`
	}
	s := ts{TS: Delivery{
		WDay: time.Sunday,
		From: Hour(13),
		To:   Hour(15),
	}}
	b, e := json.Marshal(s)
	assert.NoError(t, e)
	assert.Equal(t, `{"ts":"Sunday 1PM - 3PM"}`, string(b))

	b = []byte(`{"ts":"Sunday 1PM - 3PM"}`)
	var s2 ts
	e = json.Unmarshal(b, &s2)
	assert.NoError(t, e)
	assert.Equal(t, s, s2)
}

func TestHour(t *testing.T) {
	fullClock := [...]string{
		"12AM", "1AM", "2AM", "3AM", "4AM", "5AM", "6AM", "7AM", "8AM", "9AM", "10AM", "11AM",
		"12PM", "1PM", "2PM", "3PM", "4PM", "5PM", "6PM", "7PM", "8PM", "9PM", "10PM", "11PM",
	}
	for i := uint(0); i < 24; i++ {
		s := Hour(i).String()
		assert.Equal(t, fullClock[i], s)
	}

	for i := range fullClock {
		var h Hour
		e := h.FromString([]byte(fullClock[i]))
		assert.NoError(t, e)
		assert.Equal(t, int(h), int(Hour(i)))
	}
}
