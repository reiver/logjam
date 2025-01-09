package db

import (
	"regexp"
	"strings"
)

type MetaData struct {
	Title       string
	Description string
	Image       string
}

func FetchDataForMetaTags(ctx Context, path string) *MetaData {
	ctx.Debug("{fetchDataForMetaTags} BEGIN")
	defer ctx.Debug("{fetchDataForMetaTags} END")

	// Fetch your meta data based on the path or other conditions

	var myRecord *Record

	containsAt := strings.Contains(path, "@")
	ctx.Error("Contains '@':", containsAt)
	if containsAt {
		//get host name
		re := regexp.MustCompile(`/(@\w+)/`) // Regular expression to match '@' followed by word characters
		match := re.FindStringSubmatch(path)

		hostName := ""
		if len(match) > 1 {
			hostName = match[1] // The first submatch should be '@Zaid'
			hostName = strings.TrimPrefix(hostName, "@")

		}
		ctx.Info("HostName :", hostName)
		myRecord = &Record{
			Name:        "GreatApe",
			Description: "GreatApe is Video Conferencing Application for Fediverse",
		}

	} else {
		//get room name
		parts := strings.Split(path, "/")
		roomName := parts[len(parts)-1]
		ctx.Info("RoomName :", roomName)

		record, err := fetchRecordFromPocketBase(ctx, roomName)
		if err != nil {
			ctx.Error("Error :", roomName)
			myRecord = &Record{
				Name:        "GreatApe",
				Description: "GreatApe is Video Conferencing Application for Fediverse",
			}
		} else {
			ctx.Info("Room Name:", record.Name, "Desc: ", record.Description, "Thumbnail: ", record.Thumbnail)
			myRecord = record
		}

	}

	return &MetaData{
		Title:       myRecord.Name,
		Description: myRecord.Description,
		Image:       getImageURL(ctx, myRecord),
	}
}
