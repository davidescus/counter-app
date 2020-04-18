package webserver

import (
	"counter-app/pkg/app"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// TODO add rate limiter handler

func accessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipAddr := r.RemoteAddr
		if ip := r.Header.Get("X-FORWARDED-FOR"); ip != "" {
			ipAddr = ip
		}
		// TODO add more details here
		log.Printf("[AccessLog] %s \"%s %s",
			ipAddr,
			r.Method,
			r.RequestURI,
		)
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
			w.Header().Set("Content-Type", "application/json")
			u, err := url.Parse(r.RequestURI)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Error: ", err)
				return
			}
			params := u.Query()
			results := app.Get(params["keywords"])
			response, err := json.Marshal(results)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Error: ", err)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(response)
		case "POST":
			w.Header().Set("Content-Type", "text/plain")
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

func swaggerInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("final handler ....")
		w.Header().Set("Content-Type", "text/json; charset=utf-8")
		http.ServeFile(w, r, "./static/swagger.json")
	})
}