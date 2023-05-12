package testingStuffs

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	binarytreesrv "github.com/sparkscience/logjam/backend/srv/binarytree"
	roommapssrv "github.com/sparkscience/logjam/backend/srv/roommaps"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

type DC struct {
	sync.Mutex
	ids map[int]bool
}

func Test_JoinLeaveConcurrently(t *testing.T) {
	roomName := "roomName"
	initMap := binarytreesrv.GetMap()
	dcs := &DC{
		ids: make(map[int]bool),
	}
	{
		err := roommapssrv.RoomMaps.Set(roomName, &initMap)
		panicIfErr(err, t)
	}
	rootMap, _ := roommapssrv.RoomMaps.Get(roomName)

	//starting with broadcaster
	brwsConn := &websocket.Conn{}
	var brSocket *binarytreesrv.MySocket
	{
		rootMap.Room.Insert(brwsConn)
		err := roommapssrv.RoomMaps.Set(roomName, rootMap.Room)
		panicIfErr(err, t)

		brSocket = rootMap.Room.Get(brwsConn).(*binarytreesrv.MySocket)
		brSocket.MetaData = make(map[string]string)
		brSocket.SetName("br")
		brSocket.MetaData["streamId"] = GenRandomString(32)
		rootMap.Room.ToggleHead(brSocket.Socket)
		rootMap.Room.ToggleCanConnect(brSocket.Socket)

		err = roommapssrv.RoomMaps.Set(roomName, rootMap.Room)
		panicIfErr(err, t)
	}

	audiencesCount := 64
	AudiencesGroup := make([]*websocket.Conn, audiencesCount)
	AudiencesParentDCCh := make([]*chan bool, audiencesCount)

	for i := 0; i < audiencesCount; i++ {
		AudiencesGroup[i] = &websocket.Conn{}
		ch := make(chan bool)
		AudiencesParentDCCh[i] = &ch
		dcs.ids[i] = false
	}

	joinWG := &sync.WaitGroup{}
	leaveWG := &sync.WaitGroup{}

	addAudience := func(idx int, wsConn *websocket.Conn) {
		firstRun := true
		for {
			rootMap.Room.Insert(wsConn)
			err := roommapssrv.RoomMaps.Set(roomName, rootMap.Room)
			panicIfErr(err, t)

			ok, _ := findBroadcaster(roomName)
			if !ok {
				panicIfErr(errors.New("can't find broadcaster"), t)
			}

			_socket := rootMap.Room.Get(wsConn)
			if _socket == nil {
				return
			}
			socket := _socket.(*binarytreesrv.MySocket)
			socket.MetaData = make(map[string]string)
			socket.MetaData["idx"] = strconv.Itoa(idx)
			socket.SetName("audience" + strconv.Itoa(idx))

			targetSocketNode, err := binarytreesrv.InsertChild(socket.Socket, rootMap.Room)
			panicIfErr(err, t)
			_ = targetSocketNode.(*binarytreesrv.MySocket)

			rootMap.Room.ToggleCanConnect(socket.Socket)
			if firstRun {
				joinWG.Done()
				firstRun = false
			}
			<-*AudiencesParentDCCh[idx]
			dcs.Lock()
			if dcs.ids[idx] {
				dcs.Unlock()
				break
			} else {
				dcs.Unlock()
				println("rejoining", idx)
			}
		}
	}

	leaveAudience := func(socket *binarytreesrv.MySocket) {
		defer leaveWG.Done()
		if socket.IsBroadcaster {
			println("br")
			return
		}
		idx, _ := strconv.Atoi(socket.MetaData["idx"])
		dcs.Lock()
		dcs.ids[idx] = true
		dcs.Unlock()
		connectedSocketsList := socket.GetConnectedSocketsList()
		for _, sock := range connectedSocketsList {
			println(sock.Name, "left")
			idx, _ := strconv.Atoi(sock.MetaData["idx"])
			go func(theIdx int) {
				time.Sleep(100 * time.Millisecond)
				*AudiencesParentDCCh[idx] <- true
			}(idx)
		}

		rootMap.Room.Delete(socket.Socket)
	}

	for i := 0; i < audiencesCount; i++ {
		joinWG.Add(1)
		go addAudience(i, AudiencesGroup[i])
	}
	joinWG.Wait()
	println("added", len(rootMap.Room.Nodes()), "audiences concurrently")

	for _, conn := range rootMap.Room.Nodes() {
		leaveWG.Add(1)
		time.Sleep(10 * time.Millisecond)
		go leaveAudience(conn.(*binarytreesrv.MySocket))
	}

	leaveWG.Wait()

	println("giving it 8 sec to finish")
	time.Sleep(8 * time.Second)
	println(len(rootMap.Room.Nodes()), "node at all is here.")
}

func panicIfErr(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenRandomString(size int) string {
	b := make([]byte, size)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func findBroadcaster(roomName string) (bool, *binarytreesrv.MySocket) {
	Map, found := roommapssrv.RoomMaps.Get(roomName)
	if !found {
		println(fmt.Sprintf("could not get map for room %q when trying to find broadcaster", roomName))
		return false, nil
	}
	broadcasterLevel := Map.Room.LevelNodes(1)
	var broadcaster *binarytreesrv.MySocket
	ok := len(broadcasterLevel) == 1
	if ok {
		broadcaster = broadcasterLevel[0].(*binarytreesrv.MySocket)
	}
	return ok, broadcaster
}
