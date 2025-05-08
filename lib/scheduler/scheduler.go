package scheduler

import (
	"fmt"
	"github.com/reiver/logjam/lib/db"
	cErrors "github.com/reiver/logjam/lib/errors"
	"github.com/reiver/logjam/lib/marshal"
	"github.com/reiver/logjam/lib/neynar"
	"github.com/reiver/logjam/lib/room"
	blueskysrv "github.com/reiver/logjam/srv/bluesky"
	dbsrv "github.com/reiver/logjam/srv/db"
	neynarsrv "github.com/reiver/logjam/srv/neynar"
	roomsrv "github.com/reiver/logjam/srv/room"
	"net/http"
	"time"
)

type Scheduler struct {
}

const (
	schedulesTbl = "scheduledMeetings"
	roomTbl      = "rooms"

	//keys
	scheduledTimeKey = "scheduledTime"
	roomUIDKey       = "roomUID"
	sentKey          = "sent"

	nameKey = "name"
)

func NewSchedulerSrv() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) sendToSocials(userId, text, schedId string) error {
	defer func(sId string) {
		err := dbsrv.Repository.Update(schedulesTbl, sId, map[string]any{
			sentKey: true,
		})
		if err != nil {
			panic(err)
		}
	}(schedId)
	haveNenarAcc, err := neynarsrv.Repository.NeynarAccountExists(userId)
	if err != nil {
		return err
	}
	if haveNenarAcc {
		return neynarsrv.Repository.CreateCast(userId, neynar.CastPayload{
			Text:      text,
			ParentURL: "",
			Embeds:    nil,
		})
	}

	haveBlueSkyAcc, err := blueskysrv.Repository.AccountExists(userId)
	if err != nil {
		return nil
	}
	if haveBlueSkyAcc {
		return blueskysrv.Repository.CreatePost(userId, text)
	}

	return nil
}

func (s *Scheduler) CreateSchedule(input CreateScheduleRequestModel) error {
	room, err := roomsrv.Repository.GetRoom(input.RoomUID)
	if err != nil {
		return err
	}
	if room == nil {
		return cErrors.NewErrorWithMsg(http.StatusNotFound, "couldnt find a room with this UID")
	}
	if room.OwnerID != input.UserId {
		return cErrors.NewError(http.StatusForbidden)
	}
	_, err = dbsrv.Repository.Insert(schedulesTbl, map[string]any{scheduledTimeKey: input.DateTime, roomUIDKey: input.RoomUID})
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
				roomUID := rec[roomUIDKey].(string)
				scheduledTime := rec[scheduledTimeKey].(time.Time)
				if rec[sentKey].(bool) {
					continue
				}

				rooms, err := dbsrv.Repository.GetByFilter(roomTbl, map[string]any{
					"UID": roomUID,
				})
				if err != nil {
					panic(err)
				}
				roomData := room.RoomDTO{}
				err = marshal.MapToObj(rooms[0], &roomData)
				if err != nil {
					panic(err)
				}
				roomLink := fmt.Sprintf("https://logjam.vercel.app/log/%s", roomData.UID)

				err = s.sendToSocials(roomData.OwnerID, "meeting at "+scheduledTime.GoString()+".\nlink: "+roomLink, rec["id"].(string))
				if err != nil {
					panic(err)
				}
			}
		}
	}()
}
