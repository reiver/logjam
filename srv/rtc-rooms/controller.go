package rtcroomsrv

import (
	"github.com/reiver/logjam/lib/rtc-rooms"
	"github.com/reiver/logjam/srv/goldgorilla"
	"github.com/reiver/logjam/srv/log"
	"github.com/reiver/logjam/srv/websock"
)

var Controller = rtc_rooms.NewRoomWSController(websocksrv.WebSockSrv, Repository, goldgorillasrv.Repository, logsrv.Logger)
