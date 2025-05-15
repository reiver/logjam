package scheduler

import (
	"github.com/reiver/logjam/lib/db"
	cErrors "github.com/reiver/logjam/lib/errors"
	"github.com/reiver/logjam/lib/neynar"
	blueskysrv "github.com/reiver/logjam/srv/bluesky"
	dbsrv "github.com/reiver/logjam/srv/db"
	neynarsrv "github.com/reiver/logjam/srv/neynar"
	"net/http"
	"strconv"
	"time"
)

type Scheduler struct {
}

const (
	schedulesTbl = "scheduledMeetings"

	//keys
	dateTimeKey  = "dateTime"
	textKey      = "text"
	socialAccKey = "socialAccountId"
	sentKey      = "sent"
)

func NewSchedulerSrv() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) sendToSocials(socialAccountId, text, schedId string) error {
	defer func(sId string) {
		err := dbsrv.Repository.Update(schedulesTbl, sId, map[string]any{
			sentKey: true,
		})
		if err != nil {
			panic(err)
		}
	}(schedId)
	fid, err := strconv.ParseInt(socialAccountId, 10, 64)
	if err != nil {
		haveNenarAcc, err := neynarsrv.Repository.NeynarAccountExists(fid)
		if err != nil {
			return err
		}
		if haveNenarAcc {
			return neynarsrv.Repository.CreateCast(fid, neynar.CastPayload{
				Text:      text,
				ParentURL: "",
				Embeds:    nil,
			}, "")
		}
	}

	haveBlueSkyAcc, err := blueskysrv.Repository.AccountExists(socialAccountId)
	if err != nil {
		return nil
	}
	if haveBlueSkyAcc {
		_, err = blueskysrv.Repository.RefreshTokens(socialAccountId, "", "")
		if err != nil {
			return err
		}
		return blueskysrv.Repository.CreatePost(socialAccountId, text)
	}

	return nil
}

func (s *Scheduler) CreateSchedule(input CreateScheduleRequestModel) error {
	if input.SCAccessKeys == nil {
		return cErrors.NewErrorWithMsg(http.StatusForbidden, "social account access keys are required")
	}
	accessOk := false
	if suuid, isNeynar := input.SCAccessKeys["signerUUID"]; isNeynar {
		ok, err := neynarsrv.Repository.VerifySigner(suuid.(string))
		if err != nil {
			return err
		}
		if !ok {
			return cErrors.NewErrorWithMsg(http.StatusUnauthorized, "invalid signer uuid for "+input.SocialAccountId)
		}
		accessOk = true
	}
	if did, isNeynar := input.SCAccessKeys["did"]; isNeynar {
		at, atThere := input.SCAccessKeys["accessJwt"]
		rt, rtThere := input.SCAccessKeys["refreshJwt"]
		if !atThere || !rtThere || (len(at.(string)) == 0 || len(rt.(string)) == 0) {
			return cErrors.NewErrorWithMsg(http.StatusUnauthorized, "accessJwt and refreshJwt are required in social account access keys")
		}
		_, err := blueskysrv.Repository.RefreshTokens(did.(string), at.(string), rt.(string))
		if err != nil {
			return err
		}
		accessOk = true
	}
	if !accessOk {
		return cErrors.NewErrorWithMsg(http.StatusForbidden, "couldnt verify user access to social account with this id :"+input.SocialAccountId)
	}

	_, err := dbsrv.Repository.Insert(schedulesTbl, map[string]any{
		dateTimeKey: input.DateTime, textKey: input.Text,
		socialAccKey: input.SocialAccountId,
	})
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
		dateTimeKey + "[gte]": nowFormatted,
		dateTimeKey + "[lte]": tenMinutesFromNowFormatted,
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
				if rec[sentKey].(bool) {
					continue
				}
				err = s.sendToSocials(rec[socialAccKey].(string), rec[textKey].(string), rec["id"].(string))
				if err != nil {
					panic(err)
				}
			}
		}
	}()
}
