package websocksrv

import (
	"github.com/reiver/logjam/lib/websock"
	"github.com/reiver/logjam/srv/log"
)

var WebSockSrv websock.SocketService = websock.NewSocketService(logsrv.Logger)
