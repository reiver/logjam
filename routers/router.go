package routers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/reiver/logjam/controllers"
	"github.com/reiver/logjam/lib/logs"
	"github.com/reiver/logjam/models/contracts"
)

type IRouteRegistrar interface {
	registerRoutes(router *mux.Router)
}

type Router struct {
	router            *mux.Router
	roomWSRouter      IRouteRegistrar
	GoldGorillaRouter IRouteRegistrar
	logger            logs.ILogger
}

func NewRouter(roomWSCtrl *controllers.RoomWSController, GoldGorillaCtrl *controllers.GoldGorillaController, roomRepo contracts.IRoomRepository, socketSVC contracts.ISocketService, logger logs.ILogger) *Router {
	return &Router{
		router:            mux.NewRouter(),
		roomWSRouter:      newRoomWSRouter(roomWSCtrl, roomRepo, socketSVC, logger),
		GoldGorillaRouter: newGoldGorillaRouter(GoldGorillaCtrl),
		logger:            logger,
	}
}

// Define structs to match the JSON structure
type Record struct {
	CollectionID   string `json:"collectionId"`
	CollectionName string `json:"collectionName"`
	Created        string `json:"created"`
	Description    string `json:"description"`
	HostID         string `json:"hostId"`
	ID             string `json:"id"`
	Name           string `json:"name"`
	Thumbnail      string `json:"thumbnail"`
	Updated        string `json:"updated"`
}

type Response struct {
	Page       int      `json:"page"`
	PerPage    int      `json:"perPage"`
	TotalItems int      `json:"totalItems"`
	TotalPages int      `json:"totalPages"`
	Items      []Record `json:"items"`
}

func fetchRecordFromPocketBase(roomName string) (*Record, error) {
	url := "https://pb.greatape.stream/api/collections/rooms/records?sort=-created"

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	fmt.Println("Request:", req)

	// Send the request via a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	// fmt.Println("Response body:", body)
	record, err := parseResponse(string(body), roomName)

	return record, err
	// fmt.Println("Response:", string(body))
}

func parseResponse(jsonData string, roomName string) (*Record, error) {
	var response Response
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	for _, item := range response.Items {
		if item.Name == roomName {
			fmt.Printf("Found record: %+v\n", item)
			return &item, nil
		}
	}

	return nil, fmt.Errorf("record not found")
}

func (r *Router) RegisterRoutes() error {
	r.roomWSRouter.registerRoutes(r.router)
	r.GoldGorillaRouter.registerRoutes(r.router)

	r.router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./web-app/dist/assets/"))))
	// r.router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./web-app/dist/index.html")
	// })

	r.router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		r.logger.Log("URL: ", logs.LDebug, req.URL.Path)
		metaData := fetchDataForMetaTags(req.URL.Path) // Implement this to fetch meta data based on the request

		// Read the existing index.html file
		htmlContent, err := ioutil.ReadFile("./web-app/dist/index.html")
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			return
		}

		// Modify the HTML content to include the dynamic title and description
		modifiedHTML := injectMetaTags(string(htmlContent), metaData, r)
		// r.logger.Log("Modified HTML", logs.LDebug, modifiedHTML)

		// Serve the modified HTML content
		// w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(modifiedHTML))

		// Check if the request is from a bot
		// if isBotRequest(req.UserAgent()) {
		// 	metaData := fetchDataForMetaTags(req.URL.Path) // Implement this to fetch meta data based on the request

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
	})

	r.router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ip := request.RemoteAddr
			if strings.Index(ip, ":") > 0 {
				ip = ip[:strings.Index(ip, ":")]
			}
			r.logger.Log("HTTP", logs.LDebug, fmt.Sprintf(`%s | %s "%s"`, ip, request.Method, request.URL.Path))
			handler.ServeHTTP(writer, request)
		})
	})

	return nil
}

func isBotRequest(userAgent string) bool {
	// Simplified user agent check; extend this as needed
	lowerAgent := strings.ToLower(userAgent)
	return strings.Contains(lowerAgent, "googlebot") || strings.Contains(lowerAgent, "bingbot")
}

func fetchDataForMetaTags(path string) *MetaData {
	// Fetch your meta data based on the path or other conditions

	var myRecord *Record

	containsAt := strings.Contains(path, "@")
	fmt.Println("Contains '@':", containsAt)
	if containsAt {
		//get host name
		re := regexp.MustCompile(`/(@\w+)/`) // Regular expression to match '@' followed by word characters
		match := re.FindStringSubmatch(path)

		hostName := ""
		if len(match) > 1 {
			hostName = match[1] // The first submatch should be '@Zaid'
			hostName = strings.TrimPrefix(hostName, "@")

		}
		fmt.Println("HostName :", hostName)
		myRecord = &Record{
			Name:        "GreatApe",
			Description: "GreatApe is Video Conferencing Application for Fediverse",
		}

	} else {
		//get room name
		parts := strings.Split(path, "/")
		roomName := parts[len(parts)-1]
		fmt.Println("RoomName :", roomName)

		record, err := fetchRecordFromPocketBase(roomName)
		if err != nil {
			fmt.Println("Error :", roomName)
			myRecord = &Record{
				Name:        "GreatApe",
				Description: "GreatApe is Video Conferencing Application for Fediverse",
			}
		} else {
			fmt.Println("Room Name:", record.Name, "Desc: ", record.Description, "Thumbnail: ", record.Thumbnail)
			myRecord = record
		}

	}

	return &MetaData{
		Title:       myRecord.Name,
		Description: myRecord.Description,
		Image:       getImageURL(myRecord),
	}
}

func getImageURL(myRecord *Record) string {
	if myRecord.Thumbnail != "" {
		return "https://pb.greatape.stream/api/files/" + myRecord.CollectionID + "/" + myRecord.ID + "/" + myRecord.Thumbnail
	}
	return "" // Return an empty string or a default image URL if no thumbnail is available
}

func injectMetaTags(htmlContent string, data *MetaData, r *Router) string {
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
			r.logger.Log("IMAGE TAG", logs.LDebug, "FOUND THE IMAGE TAG")
			htmlContent = strings.Replace(htmlContent, `<meta name="twitter:image" content="/assets/metatagsLogo-3d1cffd4.png" />`, imageTag, 1)
			// r.logger.Log("UPDATED HTML: ", logs.LDebug, htmlContent)
		}

		if strings.Contains(htmlContent, `<meta property="og:image" content="/assets/metatagsLogo-3d1cffd4.png" />`) {
			r.logger.Log("IMAGE TAG", logs.LDebug, "FOUND THE IMAGE TAG")
			htmlContent = strings.Replace(htmlContent, `<meta property="og:image" content="/assets/metatagsLogo-3d1cffd4.png" />`, imageTag, 1)
			// r.logger.Log("UPDATED HTML: ", logs.LDebug, htmlContent)
		}
	}

	r.logger.Log("DESC TAG CONTENT: ", logs.LDebug, descTag)

	// Replace existing meta description tag, or add if not present
	// if strings.Contains(htmlContent, `meta property="og:title"`) {
	// 	htmlContent = strings.Replace(htmlContent, `<meta property="og:title" content="GreatApe" />`, titleTag, 1)
	// 	r.logger.Log("String found", logs.LDebug)
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
		r.logger.Log("Title and description tags inserted", logs.LDebug)
	} else {
		r.logger.Log("No <head> tag found, cannot insert tags", logs.LDebug)
	}

	return htmlContent
}

type MetaData struct {
	Title       string
	Description string
	Image       string
}

func (r *Router) Serve(addr string) error {
	println("[HTTP] serving on", addr)
	return http.ListenAndServe(addr, r.router)
}
