package httpsrv

import (
	"context"
	"github.com/reiver/logjam/lib/rest"
	"github.com/reiver/logjam/lib/tokens"
	userssrv "github.com/reiver/logjam/srv/users"
	"net/http"

	"github.com/gorilla/mux"
)

var Router *mux.Router = mux.NewRouter()
var RouterWithAuth *mux.Router

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // or your domain
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func init() {
	Router = mux.NewRouter()
	if nil == Router {
		panic("nil HTTP gorilla mux router")
	}

	//Router.Use(mux.CORSMethodMiddleware(Router))
	Router.Use(corsMiddleware)

	Router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			remoteAddr := request.RemoteAddr
			log.Debugf(`HTTP-request method=%q remote-addr=%q request-uri=%q`, request.Method, remoteAddr, request.URL.RequestURI())
			handler.ServeHTTP(writer, request)
		})
	})

	RouterWithAuth = Router.NewRoute().Subrouter()
	RouterWithAuth.Use(func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			token := r.Header.Get("Authorization")
			claims, err := tokens.ParseToken(token)
			if rest.HandleIfErr(w, err, http.StatusUnauthorized) {
				return
			}
			user, err := userssrv.Repository.GetById(claims["userId"].(string))
			if rest.HandleIfErr(w, err, http.StatusUnauthorized) {
				return
			}
			ctx := context.WithValue(r.Context(), rest.UserDataCtxKey, user)
			nextHandler.ServeHTTP(w, r.WithContext(ctx))
			//nextHandler.ServeHTTP(w, r)
		})
	})
}
