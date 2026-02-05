package websocksrv

import (
	"codeberg.org/greatape/logjam/lib/websock"
	"codeberg.org/greatape/logjam/srv/log"
)

var WebSockSrv websock.SocketService = websock.NewSocketService(logsrv.Logger)
