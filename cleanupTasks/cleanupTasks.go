package cleanupTasks

import (
	"log"
	"time"
)

type CleanupTasks []*CleanupTask

type CleanupTask struct {
	Directory string `json:"directory"`
	Date      string `json:"date"`
}

func (z *CleanupTask) GetTime(loc *time.Location) (t time.Time, err error) {

	t, err = time.ParseInLocation("2006-01-02", z.Date, loc)
	if err != nil {
		log.Fatalf("Could not parse date (%s)\n", z.Date)
	}

	return
}
