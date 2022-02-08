package graphapihandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/mmcomp/go-binarytree"
	logger "github.com/mmcomp/go-log"
)

var (
	Map = binarytree.Default
)

type MyNode struct {
	Name           string
	ID             float64
	ConnectedNodes map[float64]MyNode
	Head           bool
	Connectable    bool
}

func (receiver *MyNode) Insert(node binarytree.SingleNode) {
	theNode := node.(*MyNode)
	receiver.ConnectedNodes[theNode.ID] = *theNode
}

func (receiver *MyNode) Delete(nodeId interface{}) {
	delete(receiver.ConnectedNodes, nodeId.(float64))
}

func (receiver *MyNode) Get(nodeId interface{}) binarytree.SingleNode {
	result := receiver.ConnectedNodes[nodeId.(float64)]
	return &result
}

func (receiver *MyNode) GetLength() int {
	return len(receiver.ConnectedNodes)
}

func (receiver *MyNode) IsHead() bool {
	return receiver.Head
}

func (receiver *MyNode) CanConnect() bool {
	return receiver.Connectable
}

func (receiver *MyNode) GetAll() map[interface{}]binarytree.SingleNode {
	var output map[interface{}]binarytree.SingleNode = make(map[interface{}]binarytree.SingleNode)
	for indx := range receiver.ConnectedNodes {
		result := receiver.ConnectedNodes[indx]
		output[indx] = &result
	}
	return output
}

func (receiver *MyNode) ToggleHead() {
	receiver.Head = !receiver.Head
}

func (receiver *MyNode) ToggleCanConnect() {
	receiver.Connectable = !receiver.Connectable
}

func (receiver *MyNode) GetIndex() interface{} {
	return receiver.ID
}

func fillFunction(node interface{}, socketIndex uint64) binarytree.SingleNode {
	id := node.(float64)
	result := MyNode{
		ID:             id,
		Name:           "Node " + strconv.FormatFloat(id, 'f', 0, 64),
		ConnectedNodes: make(map[float64]MyNode),
	}
	return &result
}

func (receiver httpHandler) GetMap() binarytree.Tree {
	Map.SetFillNode(fillFunction)
	return Map
}

type TreeGraphElement struct {
	Name     string             `json:"name"`
	Parent   string             `json:"parent"`
	Children []TreeGraphElement `json:"children"`
}

func (receiver httpHandler) addSubSockets(node MyNode, children *[]TreeGraphElement, aMap binarytree.Tree) {
	for child := range node.ConnectedNodes {
		childSocket := aMap.Get(child).(*MyNode)
		*children = append(*children, TreeGraphElement{
			Name:     childSocket.Name,
			Parent:   "null",
			Children: []TreeGraphElement{},
		})
		receiver.addSubSockets(*childSocket, &(*children)[len(*children)-1].Children, aMap)
	}
}

func (receiver httpHandler) GetTree(aMap binarytree.Tree) []TreeGraphElement {
	log := receiver.Logger.Begin()
	defer log.End()

	treeData := []TreeGraphElement{}
	broadcasterLevel := aMap.LevelNodes(1)
	log.Inform("broadcasterLevel ", broadcasterLevel)
	if len(broadcasterLevel) == 0 {
		return treeData
	}
	broadcaster := broadcasterLevel[0].(*MyNode)
	treeData = append(treeData, TreeGraphElement{
		Name:     broadcaster.Name,
		Parent:   "null",
		Children: []TreeGraphElement{},
	})
	receiver.addSubSockets(*broadcaster, &treeData[0].Children, aMap)
	return treeData
}

type httpHandler struct {
	Logger logger.Logger
}

func Handler(logger logger.Logger) http.Handler {
	return httpHandler{
		Logger: logger,
	}
}

func (receiver httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log := receiver.Logger.Begin()
	defer log.End()
	Map = receiver.GetMap()

	log.Inform("Command : ", req.Method)
	w.Header().Set("Content-Type", "application/json")
	switch req.Method {
	case "GET":
		nodes := Map.GetAll()
		log.Inform("total ", len(nodes))
		treeData := receiver.GetTree(Map)
		log.Inform("treeData ", treeData)
		j, e := json.Marshal(treeData)
		if e != nil {
			log.Error(e)
			w.Write([]byte("json marshal error"))
			return
		}
		w.Write(j)
	case "POST":
		body, e := ioutil.ReadAll(req.Body)
		if e != nil {
			log.Error(e)
			w.Write([]byte("body read error"))
			return
		}
		log.Inform("body ", body)
		var bodyInterface map[string]interface{}
		e = json.Unmarshal(body, &bodyInterface)
		if e != nil {
			log.Error(e)
			w.Write([]byte("body unmarshal error"))
			return
		}
		log.Inform("bodyInterface", bodyInterface)
		ID, ok := bodyInterface["id"].(float64)
		if !ok {
			w.Write([]byte("id is not correct"))
			return
		}
		Map.Insert(ID)
		Map.ToggleCanConnect(ID)
		isHeader, hok := bodyInterface["head"].(bool)
		if hok && isHeader {
			log.Inform("Set head!")
			Map.ToggleHead(ID)
		} else {
			log.Inform("insertChild ", ID)
			Map.InsertChild(ID, false)
		}
		w.Write(body)
	case "DELETE":
		ID, err := strconv.ParseFloat(strings.Replace(req.URL.Path, "/api/v2/graph/", "", 1), 64)
		if err != nil {
			log.Error(err)
			w.Write([]byte("id is not correct"))
			return
		}
		Map.Delete(ID)
		w.Write([]byte("{\"id\":" + strconv.FormatFloat(ID, 'f', 0, 64) + "}"))
	}

}
