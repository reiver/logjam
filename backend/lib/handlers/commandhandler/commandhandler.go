package commandhandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

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

type TreeGraphElement struct {
	Name     string             `json:"name"`
	Parent   string             `json:"parent"`
	Children []TreeGraphElement `json:"children"`
}

func (receiver httpHandler) AddSubSockets(socket websocketmap.MySocket, children *[]TreeGraphElement) {
	for child := range socket.ConnectedSockets {
		childSocket := websocketmap.Map.Get(child)
		*children = append(*children, TreeGraphElement{
			Name:     childSocket.Name,
			Parent:   "null",
			Children: []TreeGraphElement{},
		})
		receiver.AddSubSockets(childSocket, &(*children)[len(*children)-1].Children)
	}
}

func (receiver httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log := receiver.Logger.Begin()
	defer log.End()

	log.Inform("Command : ", req.URL.Path)
	switch req.URL.Path {
	case "/v1/connections/graph":
		w.Header().Add("Content-Type", "text/html; charset=UTF-8")
		treeData := []TreeGraphElement{}
		var output string = ""
		fileBuffer, err := ioutil.ReadFile("./lib/handlers/commandhandler/template.html")
		if err != nil {
			log.Error("Erro reading template ", err)
			w.Write([]byte("Error reading template"))
			return
		}
		fileString := string(fileBuffer)
		brodcaster, found := websocketmap.Map.GetBroadcaster()
		if found {
			log.Alert("Br Id ", brodcaster.ID)
			treeData = append(treeData, TreeGraphElement{
				Name:     brodcaster.Name,
				Parent:   "null",
				Children: []TreeGraphElement{},
			})
			receiver.AddSubSockets(brodcaster, &treeData[0].Children)
		}
		j, e := json.Marshal(treeData)
		if e != nil {
			log.Error(e)
			w.Write([]byte("Error marshal"))
			return
		}
		output = string(j)
		fileString = strings.Replace(fileString, "#treeData#", output, 1)
		w.Write([]byte(fileString))
		return
	}
	websocketmap.Map.Reset()
	w.Write([]byte("Command '" + req.URL.Path + "' executed"))
}
