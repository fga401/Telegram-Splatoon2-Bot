package util

import "time"

// Time is a collection of functions about time
var Time = timeHelper{
	updateInterval: int64(2 * time.Hour.Seconds()),
}

type timeHelper struct {
	updateInterval int64
}

// LocalTime returns timestamp + offset with type time.Time.
func (timeHelper) LocalTime(timestamp int64, offsetInMinute int) time.Time {
	return time.Unix(timestamp, 0).In(time.FixedZone("", offsetInMinute*60))
}

// SplatoonNextUpdateTime returns the next Splatoon stage update time after given time.
func (helper timeHelper) SplatoonNextUpdateTime(after time.Time) time.Time {
	nowTimestamp := after.Unix()
	nextTimestamp := (nowTimestamp/helper.updateInterval + 1) * helper.updateInterval
	return time.Unix(nextTimestamp, 0)
}
