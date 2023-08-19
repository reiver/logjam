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
	//Children [2]MemberModel
	Children []MemberModel
}

type PeerModel struct {
	ID              uint64
	IsConnected     bool
	IsAuxiliaryNode bool
	//Children        [2]*PeerModel
	Children []*PeerModel
}

type RoomModel struct {
	*sync.Mutex
	Title                  string
	PeersTree              *PeerModel
	Members                map[uint64]*MemberModel
	MetaData               map[string]any
	AuxiliaryNode          **PeerModel
	HadAuxiliaryNodeBefore bool
}

func (r *RoomModel) GetBroadcaster() *MemberModel {
	if r.PeersTree.IsConnected {
		return r.Members[r.PeersTree.ID]
	}
	return nil
}

func (r *RoomModel) GetLevelMembers(level uint, includeAll bool) ([]**PeerModel, error) {
	if level == 0 {
		if !includeAll && r.PeersTree.IsConnected {
			return []**PeerModel{&r.PeersTree}, nil
		} else if includeAll {
			return []**PeerModel{&r.PeersTree}, nil
		}
		return []**PeerModel{}, nil
	}
	var lastLevelChildren []**PeerModel

	if len(r.PeersTree.Children) > 0 {
		if !includeAll && r.Members[r.PeersTree.Children[0].ID].CanAcceptChild {
			lastLevelChildren = append(lastLevelChildren, &r.PeersTree.Children[0])
		} else if includeAll {
			lastLevelChildren = append(lastLevelChildren, &r.PeersTree.Children[0])
		}
	}
	if len(r.PeersTree.Children) > 1 {
		if !includeAll && r.Members[r.PeersTree.Children[1].ID].CanAcceptChild {
			lastLevelChildren = append(lastLevelChildren, &r.PeersTree.Children[1])
		} else if includeAll {
			lastLevelChildren = append(lastLevelChildren, &r.PeersTree.Children[1])
		}
	}
	if level == 1 {
		return lastLevelChildren, nil
	}
	currentLevel := uint(1)
	for currentLevel < level && len(lastLevelChildren) > 0 {
		var newList []**PeerModel
		for _, currentLevelChild := range lastLevelChildren {
			for _, child := range (*currentLevelChild).Children {
				if child == nil {
					continue
				}
				if !includeAll && r.Members[child.ID].CanAcceptChild {
					newList = append(newList, &child)
				} else if includeAll {
					newList = append(newList, &child)
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