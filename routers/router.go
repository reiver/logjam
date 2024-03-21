package routers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"sourcecode.social/greatape/logjam/controllers"
	"sourcecode.social/greatape/logjam/models/contracts"
)

type IRouteRegistrar interface {
	registerRoutes(router *mux.Router)
}

type Router struct {
	router            *mux.Router
	roomWSRouter      IRouteRegistrar
	GoldGorillaRouter IRouteRegistrar
	logger            contracts.ILogger
}

func NewRouter(roomWSCtrl *controllers.RoomWSController, GoldGorillaCtrl *controllers.GoldGorillaController, roomRepo contracts.IRoomRepository, socketSVC contracts.ISocketService, logger contracts.ILogger) *Router {
	return &Router{
		router:            mux.NewRouter(),
		roomWSRouter:      newRoomWSRouter(roomWSCtrl, roomRepo, socketSVC, logger),
		GoldGorillaRouter: newGoldGorillaRouter(GoldGorillaCtrl),
		logger:            logger,
	}
}

func (r *Router) RegisterRoutes() error {
	r.roomWSRouter.registerRoutes(r.router)
	r.GoldGorillaRouter.registerRoutes(r.router)

	r.router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./web-app/dist/assets/"))))
	// r.router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./web-app/dist/index.html")
	// })

	r.router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		metaData := fetchDataForMetaTags(req.URL.Path) // Implement this to fetch meta data based on the request

		// Read the existing index.html file
		htmlContent, err := ioutil.ReadFile("./web-app/dist/index.html")
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			return
		}

		// Modify the HTML content to include the dynamic title and description
		modifiedHTML := injectMetaTags(string(htmlContent), metaData, r)
		r.logger.Log("Modified HTML", contracts.LDebug, modifiedHTML)

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
			r.logger.Log("HTTP", contracts.LDebug, fmt.Sprintf(`%s | %s "%s"`, ip, request.Method, request.URL.Path))
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
	return &MetaData{
		Title:       "Dynamic Title",
		Description: "Dynamic Description",
	}
}

func injectMetaTags(htmlContent string, data *MetaData, r *Router) string {
	// Inject title and meta description into the HTML content
	// descriptionTag := `<meta property="og:description" content="` + data.Description + `">`
	titleTag := `<meta property="og:title" content="` + data.Title + `">`

	r.logger.Log("HTML CONTENT", contracts.LDebug, htmlContent)

	// Replace existing meta description tag, or add if not present
	if strings.Contains(htmlContent, `meta property="og:title"`) {
		htmlContent = strings.Replace(htmlContent, `<meta property="og:title" content="GreatApe" />`, titleTag, 1)
		r.logger.Log("String found", contracts.LDebug)
	}
	// else {
	// 	htmlContent = strings.Replace(htmlContent, "<head>", "<head>"+descriptionTag, 1)
	// 	r.logger.Log("String Not found", contracts.LDebug)

	// }

	return htmlContent
}

type MetaData struct {
	Title       string
	Description string
}

func (r *Router) Serve(addr string) error {
	println("[HTTP] serving on", addr)
	return http.ListenAndServe(addr, r.router)
}
