package roomRepository

import (
	"sourcecode.social/greatape/logjam/helpers"
	"testing"
)

func Test_CreationAndExistence(t *testing.T) {
	repo := NewRoomRepository()
	roomId := "roomId"
	memberId := uint64(2)
	memberEmail := "mail@example.com"
	memberName := "memberName"
	memberStreamId := "abcdefghijklmnopqrstuvwxyz1234567890"
	err := repo.CreateRoom(roomId)
	if err != nil {
		t.Error(err)
		return
	}
	if !repo.DoesRoomExists(roomId) {
		t.Error("room doesn't exists after creation")
	}
	err = repo.AddMember(roomId, memberId, memberName, memberEmail, memberStreamId, false)
	if err != nil {
		t.Error(err)
		return
	}
	member, err := repo.GetMember(roomId, memberId)
	if err != nil {
		t.Error(err)
		return
	}
	if member == nil {
		t.Error("member doesn't exists after adding to room")
		return
	}
	if member.Name != memberName ||
		member.ID != memberId ||
		member.Email != member.Email ||
		member.MetaData["streamId"] != memberStreamId {
		t.Error("member info doesn't match")
		return
	}
	membersIdList, err := repo.GetAllMembersId(roomId, false)
	if err != nil {
		t.Error(err)
		return
	}
	if len(membersIdList) != 1 || membersIdList[0] != memberId {
		t.Error("[repo.GetAllMembersId()] failed, list length doesnt match or member is not in the list")
		return
	}
	membersList, err := repo.GetMembersList(roomId)
	if err != nil {
		t.Error(err)
		return
	}
	if len(membersList) != 1 || membersList[0].Id != memberId {
		t.Error("[repo.GetMembersList()] failed, list length doesnt match or member is not in the list")
		return
	}
	memberName = "newName"
	err = repo.UpdateMemberName(roomId, memberId, memberName)
	if err != nil {
		t.Error(err)
		return
	}
	member, err = repo.GetMember(roomId, memberId)
	if err != nil {
		t.Error(err)
		return
	}
	if member.Name != memberName {
		t.Error("[repo.UpdateMemberName()] didnt updated the member name correctly or GetMember() returned wrong info")
		return
	}
	err = repo.UpdateMemberMeta(roomId, memberId, "testyKey", "testyValue")
	if err != nil {
		t.Error(err)
		return
	}
	if testykv, exists := member.MetaData["testyKey"]; exists {
		if testykv != "testyValue" {
			t.Error("[repo.UpdateMemberMeta()]wrong metadata value")
		}
	} else {
		t.Error("metadata is not updated or returned member from GetMember() func is not referenced correctly")
		return
	}
}

func Test_Tree(t *testing.T) {
	repo := NewRoomRepository()
	roomId := "roomId"
	err := repo.CreateRoom(roomId)
	if err != nil {
		t.Error(err)
		return
	}
	err = repo.AddMember(roomId, uint64(0), helpers.GetRandomString(8), "", "", false)
	if err != nil {
		t.Error(err)
		return
	}
	err = repo.UpdateCanConnect(roomId, uint64(0), true)
	if err != nil {
		t.Error(err)
		return
	}
	err = repo.SetBroadcaster(roomId, uint64(0))
	if err != nil {
		t.Error(err)
		return
	}

	err = repo.AddMember(roomId, uint64(1), helpers.GetRandomString(8), "", "", false)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = repo.InsertMemberToTree(roomId, uint64(1), false)
	if err != nil {
		t.Error(err)
		return
	}
	err = repo.AddMember(roomId, uint64(2), helpers.GetRandomString(8), "", "", false)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = repo.InsertMemberToTree(roomId, uint64(2), false)
	if err != nil {
		t.Error(err)
		return
	}
	err = repo.AddMember(roomId, uint64(3), helpers.GetRandomString(8), "", "", false)
	if err != nil {
		t.Error(err)
		return
	}
	err = repo.UpdateCanConnect(roomId, uint64(3), true)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = repo.InsertMemberToTree(roomId, uint64(3), false)
	if err == nil {
		t.Error("shouldn't be added to tree as there is no active node to connect to")
		return
	}
	err = repo.UpdateCanConnect(roomId, uint64(1), true)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = repo.InsertMemberToTree(roomId, uint64(3), false)
	if err != nil {
		t.Error(err)
		return
	}
	err = repo.UpdateCanConnect(roomId, uint64(2), true)
	if err != nil {
		t.Error(err)
		return
	}
	room, err := repo.GetRoom(roomId)
	if err != nil {
		t.Error(err)
		return
	}
	levelNodes, err := room.GetLevelMembers(0, false)
	if err != nil {
		t.Error(err)
		return
	}
	if len(levelNodes) != 1 || (*levelNodes[0]).ID != 0 || !(*levelNodes[0]).IsConnected {
		t.Error("broadcaster is not at the seat!")
	}
	levelNodes, err = room.GetLevelMembers(1, false)
	if err != nil {
		t.Error(err)
		return
	}
	if len(levelNodes) != 2 || ((*levelNodes[0]).ID != 1 && (*levelNodes[0]).ID != 2) {
		t.Error("can't find two level 1 nodes on the tree")
	}
	levelNodes, err = room.GetLevelMembers(2, false)
	if err != nil {
		t.Error(err)
		return
	}
	if len(levelNodes) != 1 || (*levelNodes[0]).ID != 3 {
		t.Error("can't find target level 2 audience")
	}
	wasBroadcaster, children, err := repo.RemoveMember(roomId, 3)
	if err != nil {
		t.Error(err)
		return
	}
	if wasBroadcaster {
		t.Error("removed member wasn't broadcaster but repo.RemoveMember() results opposite")
		return
	}
	if len(children) > 0 {
		t.Error("this node had no children but repo.RemoveMember() results with children list")
		return
	}
	levelNodes, err = room.GetLevelMembers(2, true)
	if err != nil {
		t.Error(err)
		return
	}
	if len(levelNodes) > 0 {
		t.Error("there shouldn't be any level 2 member in the tree")
	}
}
