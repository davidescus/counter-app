package publicapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// TODO add rate limiter handler

func accessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipAddr := r.RemoteAddr
		if ip := r.Header.Get("X-FORWARDED-FOR"); ip != "" {
			ipAddr = ip
		}
		// TODO add more details here
		log.Printf("[AccessLog] %s \"%s %s %s",
			ipAddr,
			r.Method,
			r.RequestURI,
			r.UserAgent(),
		)
		next.ServeHTTP(w, r)
	})
}

// use it if needed
func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func final(storage Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json")
			u, err := url.Parse(r.RequestURI)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("[Error] ", err)
				return
			}
			params := u.Query()
			results := get(params["keywords"], storage)
			response, err := json.Marshal(results)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("[Error] ", err)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(response)
		case "POST":
			w.Header().Set("Content-Type", "text/plain")
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("[Error] ", err)
				return
			}
			store(body, storage)
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})
}

func swaggerInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/json; charset=utf-8")
		http.ServeFile(w, r, "./static/swagger.json")
	})
}

func get(keywords []string, storage Storage) []Occurrence {
	var results []Occurrence
	for _, key := range keywords {
		results = append(results, Occurrence{
			Key:   key,
			Count: storage.Get([]byte(strings.ToLower(key))),
		})
	}
	return results
}

func store(body []byte, storage Storage) {
	keywords := strings.Fields(strings.ToLower(string(body)))
	for _, keyword := range keywords {
		storage.Increment([]byte(keyword))
	}
}
