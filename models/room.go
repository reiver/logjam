package models

import (
	"sync"
)

type MemberModel struct {
	ID             uint64
	Name           string
	Email          string
	IsUsingTurn    bool
	MetaData       map[string]any
	CanAcceptChild bool
}

type PeerTreeModel struct {
	MemberModel
	Children [2]MemberModel
}

type PeerModel struct {
	ID          uint64
	IsConnected bool
	Children    [2]*PeerModel
}

type RoomModel struct {
	*sync.Mutex
	Title     string
	PeersTree *PeerModel
	Members   map[uint64]*MemberModel
	MetaData  map[string]any
}

func (r *RoomModel) GetBroadcaster() *MemberModel {
	if r.PeersTree.IsConnected {
		return r.Members[r.PeersTree.ID]
	}
	return nil
}

func (r *RoomModel) GetLevelAudiences(level uint) ([]*PeerModel, error) {
	if level == 0 {
		if r.PeersTree.IsConnected {
			return []*PeerModel{r.PeersTree}, nil
		}
		return []*PeerModel{}, nil
	}
	var lastLevelChildren = r.PeersTree.Children[:]
	if lastLevelChildren[0] == nil || !r.Members[lastLevelChildren[0].ID].CanAcceptChild {
		lastLevelChildren = append(lastLevelChildren[:0], lastLevelChildren[1:]...)
	}
	if lastLevelChildren[0] == nil || !r.Members[lastLevelChildren[0].ID].CanAcceptChild {
		lastLevelChildren = append(lastLevelChildren[:0], lastLevelChildren[1:]...)
	}
	if level == 1 {
		return lastLevelChildren, nil
	}
	currentLevel := uint(1)
	for currentLevel < level && len(lastLevelChildren) > 0 {
		var newList []*PeerModel
		for _, currentLevelChild := range lastLevelChildren {
			for _, nextLevelChild := range currentLevelChild.Children {
				if nextLevelChild != nil && r.Members[nextLevelChild.ID].CanAcceptChild {
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
	return lastLevelChildren, nil
}

func (r *RoomModel) GetTree() (*PeerTreeModel, error) {

	return nil, nil
}
