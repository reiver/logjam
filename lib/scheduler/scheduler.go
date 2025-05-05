package scheduler

import (
	"fmt"
	"github.com/reiver/logjam/lib/db"
	dbsrv "github.com/reiver/logjam/srv/db"
	"time"
)

type Scheduler struct {
}

const (
	schedulesTbl = "scheduledMeetings"
	roomTbl      = "rooms"

	//keys
	scheduledTimeKey = "scheduledTime"
	roomIdKey        = "roomId"

	nameKey = "name"
)

func NewSchedulerSrv() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) CreateSchedule(reqModel CreateScheduleRequestModel) error {
	_, err := dbsrv.Repository.Insert(schedulesTbl, map[string]any{scheduledTimeKey: reqModel.DateTime, roomIdKey: reqModel.RoomID})
	return err
}

func (s *Scheduler) GetNotificationsToSend() ([]db.Record, error) {
	now := time.Now()
	tenMinutesFromNow := now.Add(10 * time.Minute)

	// Format time in PocketBase compatible format
	nowFormatted := now.Format("2006-01-02 15:04:05")
	tenMinutesFromNowFormatted := tenMinutesFromNow.Format("2006-01-02 15:04:05")

	// Build the filter expression for `scheduledTime` between now and 10 minutes ahead
	filter := map[string]any{
		scheduledTimeKey + "[gte]": nowFormatted,
		scheduledTimeKey + "[lte]": tenMinutesFromNowFormatted,
	}

	// Use the filter to get the records
	records, err := dbsrv.Repository.GetByFilter(schedulesTbl, filter)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s *Scheduler) Start() {
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			recs, err := s.GetNotificationsToSend()
			if err != nil {
				panic(err)
			}
			for _, rec := range recs {
				roomId := rec[roomIdKey].(string)
				scheduledTime := rec[scheduledTimeKey].(time.Time)

				room, err := dbsrv.Repository.GetById(roomTbl, roomId)
				if err != nil {
					panic(err)
				}
				roomLink := fmt.Sprintf("https://logjam.vercel.app/log/%s", room[nameKey])

				_ = scheduledTime
				_ = roomLink
			}
		}
	}()
}
