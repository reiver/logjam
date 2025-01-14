package verboten

import (
	"net/http"
	"strings"

	"github.com/reiver/logjam/cfg"
	"github.com/reiver/logjam/lib/db"
	"github.com/reiver/logjam/srv/http"
	"github.com/reiver/logjam/srv/log"
	"github.com/reiver/logjam/web-app"
)

const pathprefix  string = "/log/"
const pathpattern string = "/@{handle}/host"

func init() {
	httpsrv.Router.PathPrefix(pathprefix).HandlerFunc(serveHTTP)
	httpsrv.Router.HandleFunc(pathpattern, serveHTTP)
}

func serveHTTP(w http.ResponseWriter, req *http.Request) {

		log.Debug("URL: ", req.URL.Path)
		var metaData *db.MetaData
		{
			var ctx = db.Context{
				Logger: logsrv.Tag("db"),
				PocketBaseURL: cfg.Config.PocketBaseURL(),
			}

			metaData = db.FetchDataForMetaTags(ctx, req.URL.Path) // Implement this to fetch meta data based on the request
		}

		// Read the existing index.html file
		htmlContent := webapp.IndexFile

		// Modify the HTML content to include the dynamic title and description
		modifiedHTML := injectMetaTags(htmlContent, metaData)
		// log.Debug("Modified HTML", modifiedHTML)

		// Serve the modified HTML content
		// w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(modifiedHTML))

		// Check if the request is from a bot
		// if isBotRequest(req.UserAgent()) {
		// 	metaData := db.FetchDataForMetaTags(req.URL.Path) // Implement this to fetch meta data based on the request

		// 	// Read the existing index.html file
		// 	htmlContent, err := ioutil.ReadFile("./web-app/dist/index.html")
		// 	if err != nil {
		// 		http.Error(w, "Internal Server Error", 500)
		// 		return
		// 	}

		// 	// Modify the HTML content to include the dynamic title and description
		// 	modifiedHTML := injectMetaTags(string(htmlContent), metaData)

		// 	// Serve the modified HTML content
		// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// 	w.Write([]byte(modifiedHTML))
		// } else {
		// 	// Serve the static index.html for regular users
		// 	http.ServeFile(w, req, "./web-app/dist/index.html")
		// }
}

func isBotRequest(userAgent string) bool {
	// Simplified user agent check; extend this as needed
	lowerAgent := strings.ToLower(userAgent)
	return strings.Contains(lowerAgent, "googlebot") || strings.Contains(lowerAgent, "bingbot")
}

func injectMetaTags(htmlContent string, data *db.MetaData) string {
	// Inject title and meta description into the HTML content

	if data.Title == "" {
		data.Title = "GreatApe"
	}

	if data.Description == "" {
		data.Description = "GreatApe is Video Conferencing Application for Fediverse"
	}

	titleTag := `<meta property="og:title" content="` + data.Title + `">`
	descTag := `<meta property="og:description" content="` + data.Description + `">`
	twitterTitleTag := `<meta property="twitter:title" content="` + data.Title + `">`
	twitterDescTag := `<meta property="twitter:description" content="` + data.Description + `">`
	imageTag := ``
	if data.Image != "" {
		imageTag = `<meta property="og:image" content="` + data.Image + `" />`

		// Remove the existing meta OG tag
		if strings.Contains(htmlContent, `<meta name="twitter:image" content="/assets/metatagsLogo-3d1cffd4.png" />`) {
			log.Debug("IMAGE TAG", "FOUND THE IMAGE TAG")
			htmlContent = strings.Replace(htmlContent, `<meta name="twitter:image" content="/assets/metatagsLogo-3d1cffd4.png" />`, imageTag, 1)
			// log.Debug("UPDATED HTML: ", htmlContent)
		}

		if strings.Contains(htmlContent, `<meta property="og:image" content="/assets/metatagsLogo-3d1cffd4.png" />`) {
			log.Debug("IMAGE TAG", "FOUND THE IMAGE TAG")
			htmlContent = strings.Replace(htmlContent, `<meta property="og:image" content="/assets/metatagsLogo-3d1cffd4.png" />`, imageTag, 1)
			// log.Debug("UPDATED HTML: ", htmlContent)
		}
	}

	log.Debug("DESC TAG CONTENT: ", descTag)

	// Replace existing meta description tag, or add if not present
	// if strings.Contains(htmlContent, `meta property="og:title"`) {
	// 	htmlContent = strings.Replace(htmlContent, `<meta property="og:title" content="GreatApe" />`, titleTag, 1)
	// 	log.Debug("String found")
	// }

	headStartIndex := strings.Index(htmlContent, "<head>")
	if headStartIndex != -1 {
		// Position to insert after <head> tag
		insertPosition := headStartIndex + len("<head>")

		// Concatenate titleTag and descTag for insertion
		tagsToInsert := titleTag + descTag + twitterTitleTag + twitterDescTag

		// if imageTag != `` {
		// 	tagsToInsert += imageTag
		// }

		// Insert the tags right after the <head> tag
		htmlContent = htmlContent[:insertPosition] + tagsToInsert + htmlContent[insertPosition:]
		log.Debug("Title and description tags inserted")
	} else {
		log.Debug("No <head> tag found, cannot insert tags")
	}

	return htmlContent
}
