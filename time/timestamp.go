package time

import "time"

// timestamp is unixnano
func TimestampToTime(timestamp int64) time.Time {
	return time.Unix(0, timestamp)
}
