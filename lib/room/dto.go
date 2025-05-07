package room

import (
	dbsrv "github.com/reiver/logjam/srv/db"
)

type CreateRoomDTO struct {
	OwnerID     string `json:"ownerId"`
	UID         string `json:"UID"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Thumbnail   string `json:"thumbnail,omitempty"`
}

type UpdateRoomDTO struct {
	OwnerID     string  `json:"-"`
	UID         string  `json:"UID" validate:"required"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Thumbnail   *string `json:"thumbnail,omitempty"`
}

type RoomDTO struct {
	ID          string        `json:"id"`
	UID         string        `json:"UID"`
	OwnerID     string        `json:"ownerId"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Thumbnail   string        `json:"thumbnail"`
	Created     dbsrv.PBTime  `json:"created"`
	Updated     *dbsrv.PBTime `json:"updated"`
}
