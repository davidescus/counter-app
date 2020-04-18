package webserver

import (
	"counter-app/pkg/app"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// TODO add rate limiter handler

func accessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO add data here, should connect request with response, to log status code for response
		log.Println("TODO access log handler add details here ...")
		next.ServeHTTP(w, r)
	})
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Use auth handler if needed ...")
		next.ServeHTTP(w, r)
	})
}

func final(app *app.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("final handler ....")
		switch r.Method {
		case "GET":
			err := r.ParseForm()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Error: ", err)
				return
			}
			keys := r.Form.Get("keywords")
			results := app.Get(strings.Fields(keys))
			response, err := json.Marshal(results)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Error: ", err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(response)
		case "POST":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Error: ", err)
				return
			}
			if err = app.Store(string(body)); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Error: ", err)
				return
			}
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})
}

func swaggerInfoHandler(app *app.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("final handler ....")
		w.Write([]byte("TODO serve static html swagger"))
	})
}
