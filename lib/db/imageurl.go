package db

func getImageURL(myRecord *Record) string {
	if myRecord.Thumbnail != "" {
		return "https://pb.greatape.stream/api/files/" + myRecord.CollectionID + "/" + myRecord.ID + "/" + myRecord.Thumbnail
	}
	return "" // Return an empty string or a default image URL if no thumbnail is available
}
