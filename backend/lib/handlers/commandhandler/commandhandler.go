package commandhandler

import (
	"net/http"

	logger "github.com/mmcomp/go-log"
	"github.com/sparkscience/logjam/backend/lib/websocketmap"
)

type httpHandler struct {
	Logger logger.Logger
}

func Handler(logger logger.Logger) http.Handler {
	return httpHandler{
		Logger: logger,
	}
}

func (receiver httpHandler) levelSockets(level uint) []websocketmap.MySocket {

	var output []websocketmap.MySocket = []websocketmap.MySocket{}
	broadCaster, ok := websocketmap.Map.GetBroadcaster()
	if ok {
		output = append(output, broadCaster)
	}
	if level == 1 || !ok {
		return output
	}
	var index uint = 1
	var currentLevelSockets []websocketmap.MySocket = output
	for {
		if len(currentLevelSockets) == 0 {
			break
		}
		output = []websocketmap.MySocket{}
		for _, socks := range currentLevelSockets {
			for _, child := range socks.ConnectedSockets {
				output = append(output, child)
			}
		}
		if len(output) == 0 {
			break
		}
		if index == level-1 {
			return output
		}
		currentLevelSockets = output
		output = []websocketmap.MySocket{}
		index++
	}
	return output
}

func (receiver httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log := receiver.Logger.Begin()
	defer log.End()

	log.Inform("Command : ", req.URL.Path)
	switch req.URL.Path {
	case "/v1/connections/graph":
		var output string = ""
		brodcaster, found := websocketmap.Map.GetBroadcaster()
		if found {
			log.Alert("Br Id ", brodcaster.ID)
			output += /*strconv.FormatUint(brodcaster.ID, 10)*/ brodcaster.Name + "<br/>\n"
			var level uint = 2
			for {
				sockets := receiver.levelSockets(level)

				if len(sockets) == 0 {
					break
				}

				for i := 0; i < len(sockets); i++ {
					socket := sockets[i]
					var hasStream string = "false"
					if socket.HasStream {
						hasStream = "true"
					}
					output += "(" + socket.Name + "[" + websocketmap.Map.GetParent(socket.Socket).Name + "] hasStream : " + hasStream + ")" //strconv.FormatUint(socket.ID, 10)
					output += ","
				}
				output += "<br/>\n"
				level++
			}
		}
		w.Header().Add("Content-Type", "text/html; charset=UTF-8")
		w.Write([]byte(output))
		return
	}
	websocketmap.Map.Reset()
	w.Write([]byte("Command '" + req.URL.Path + "' executed"))
}
