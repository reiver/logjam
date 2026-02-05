package roomsrv

import (
	"codeberg.org/greatape/logjam/lib/rooms"
	"codeberg.org/greatape/logjam/srv/goldgorilla"
	"codeberg.org/greatape/logjam/srv/log"
	"codeberg.org/greatape/logjam/srv/websock"
)

var Controller = rooms.NewRoomWSController(websocksrv.WebSockSrv, Repository, goldgorillasrv.Repository, logsrv.Logger)
