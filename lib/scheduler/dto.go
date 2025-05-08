package scheduler

import "time"

type CreateScheduleRequestModel struct {
	UserId   string    `json:"-"`
	DateTime time.Time `json:"dateTime"`
	RoomUID  string    `json:"roomUID"`
}
