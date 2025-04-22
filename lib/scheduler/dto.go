package scheduler

import "time"

type CreateScheduleRequestModel struct {
	DateTime time.Time `json:"date_time"`
	RoomID   string    `json:"room_id"`
}
