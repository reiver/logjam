package room

import "time"

type CreateRoomDTO struct {
	OwnerID     string `json:"owner_id"`
	UID         string `json:"uid"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Thumbnail   string `json:"thumbnail,omitempty"`
}

type UpdateRoomDTO struct {
	OwnerID     string  `json:"owner_id" validate:"required"`
	UID         string  `json:"uid" validate:"required"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Thumbnail   *string `json:"thumbnail,omitempty"`
}

type RoomDTO struct {
	ID          string    `json:"id"`
	UID         string    `json:"uid"`
	OwnerID     string    `json:"owner_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Thumbnail   string    `json:"thumbnail"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}
