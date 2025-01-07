package webapp

import (
	"embed"
)

///go:embed dist/assets/*.avi
///go:embed dist/assets/*.bmp
///go:embed dist/assets/*.csv
///go:embed dist/assets/*.gif
///go:embed dist/assets/*.htm
///go:embed dist/assets/*.html
///go:embed dist/assets/*.jpeg
///go:embed dist/assets/*.pdf
///go:embed dist/assets/*.txt

const AssetsPathPrefix string = "dist/assets"

//go:embed dist/assets/*.css
//go:embed dist/assets/*.ico
//go:embed dist/assets/*.jpg
//go:embed dist/assets/*.js
//go:embed dist/assets/*.png
//go:embed dist/assets/*.svg
var AssetsFS embed.FS
