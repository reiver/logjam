package scheduler

import "time"

type CreateScheduleRequestModel struct {
	DateTime        time.Time      `json:"dateTime"`
	Text            string         `json:"text"`
	SocialAccountId string         `json:"socialAccountId"`
	SCAccessKeys    map[string]any `json:"SCAccessKeys"`
}
