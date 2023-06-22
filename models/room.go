package models

import (
	"sync"
)

type UserModel struct {
	ID             uint64
	Name           string
	Email          string
	StreamId       string
	MetaData       map[string]any
	CanAcceptChild bool
}

type PeerTreeModel struct {
	ID          uint64
	IsConnected bool
	Children    [2]*PeerTreeModel
}

type RoomModel struct {
	*sync.Mutex
	Title     string
	PeersTree *PeerTreeModel
	Members   map[uint64]*UserModel
	MetaData  map[string]any
}

func (r *RoomModel) GetBroadcaster() *UserModel {
	if r.PeersTree.IsConnected {
		return r.Members[r.PeersTree.ID]
	}
	return nil
}

func (r *RoomModel) GetLevelAudiences(level uint) ([]*UserModel, error) {
	currentLevel := uint(1)
	var lastLevelChildren = r.PeersTree.Children[:]
	if lastLevelChildren[0] == nil {
		lastLevelChildren = append(lastLevelChildren[:0], lastLevelChildren[1:]...)
	}
	if lastLevelChildren[0] == nil {
		lastLevelChildren = append(lastLevelChildren[:0], lastLevelChildren[1:]...)
	}

	for currentLevel <= level && len(lastLevelChildren) > 0 {
		var newList []*PeerTreeModel
		for _, currentLevelChild := range lastLevelChildren {
			for _, nextLevelChild := range currentLevelChild.Children {
				if nextLevelChild != nil {
					newList = append(newList, nextLevelChild)
				}
			}
		}
		lastLevelChildren = newList
		if currentLevel < level {
			currentLevel++
			continue
		}
		break
	}
	var userInfoList []*UserModel
	for _, peer := range lastLevelChildren {
		userInfoList = append(userInfoList, r.Members[peer.ID])
	}
	return userInfoList, nil
}
