package worker

import (
	"time"

	c "farmer/internal/pkg/constants"
)

type workerSetting struct {
	sleepAfterQuery   time.Duration
	timeFrameUnixMili uint64
	tciLen            uint64
}

func newWorkerSetting(timeFrame string) *workerSetting {
	var (
		sleepAfterQuery   time.Duration
		timeFrameUnixMili uint64
		tciLen            uint64
	)

	switch timeFrame {
	case c.M1:
		sleepAfterQuery = time.Second
		timeFrameUnixMili = uint64(60000)
		tciLen = c.M1TciLen
	case c.H1:
		sleepAfterQuery = time.Second * 3
		timeFrameUnixMili = uint64(3600000)
		tciLen = c.H1TciLen
	default:
		// FIXME: assume not fall in this case.
		break
	}

	return &workerSetting{
		sleepAfterQuery:   sleepAfterQuery,
		timeFrameUnixMili: timeFrameUnixMili,
		tciLen:            tciLen,
	}
}
