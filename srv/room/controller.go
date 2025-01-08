package roomsrv

import (
	"github.com/reiver/logjam/lib/rooms"
	"github.com/reiver/logjam/srv/goldgorilla"
	"github.com/reiver/logjam/srv/log"
	"github.com/reiver/logjam/srv/websock"
)

var Controller = rooms.NewRoomWSController(websocksrv.WebSockSrv, Repository, goldgorillasrv.Repository, logsrv.Logger)
