package db

func getImageURL(ctx Context, myRecord *Record) string {
	if myRecord.Thumbnail != "" {
		var url string = ctx.URL("/files/" + myRecord.CollectionID + "/" + myRecord.ID + "/" + myRecord.Thumbnail)
		return url
	}
	return "" // Return an empty string or a default image URL if no thumbnail is available
}
