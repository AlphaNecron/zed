package tests

import (
	"github.com/stretchr/testify/assert"
	"maps"
	"necron.dev/zed"
	"slices"
	"testing"
	"time"
)

func TestDateTime(t *testing.T) {
	vi, e := time.LoadLocation("Asia/Ho_Chi_Minh")
	assert.NoError(t, e, "error loading location data")
	testData := map[string]time.Time{
		"2006-12-12T00:00:00+07:00": time.Date(2006, 12, 12, 0, 0, 0, 0, vi),
		"2006-11-03T00:00:00+07:00": time.Date(2006, 11, 3, 0, 0, 0, 0, vi),
	}
	actualTestData := slices.Collect(maps.Keys(testData))
	antitheses := []any{
		"foo",
		"bar",
		0,
		1,
	}
	z := zed.New().
		Field("foo", zed.DateTime("expected rfc3399 datetime value").Layout(time.RFC3339))
	testMulti[string, time.Time](t, z, "foo", actualTestData, false, func(a string, b time.Time) (time.Time, time.Time, bool) {
		return b, testData[a], testData[a].Compare(b) == 0
	})
	testMulti[any, time.Time](t, z, "foo", antitheses, true, nil)
}

func TestDateTimeEpoch(t *testing.T) {
	vi, e := time.LoadLocation("Asia/Ho_Chi_Minh")
	assert.NoError(t, e, "error loading location data")
	testData := map[float64]time.Time{
		1165856400000: time.Date(2006, 12, 12, 0, 0, 0, 0, vi),
		1162486800000: time.Date(2006, 11, 3, 0, 0, 0, 0, vi),
	}
	actualTestData := slices.Collect(maps.Keys(testData))
	antitheses := []any{
		"foo",
		"bar",
		true,
		false,
	}
	z := zed.New().
		Field("foo", zed.DateTime("expected epoch datetime value").EpochUnit(zed.EpochMillisecond))
	testMulti[float64, time.Time](t, z, "foo", actualTestData, false, func(a float64, b time.Time) (time.Time, time.Time, bool) {
		return b, testData[a], testData[a].Compare(b) == 0
	})
	testMulti[any, time.Time](t, z, "foo", antitheses, true, nil)
}
