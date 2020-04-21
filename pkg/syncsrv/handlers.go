package syncsrv

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// TODO add rate limiter handler

func accessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipAddr := r.RemoteAddr
		if ip := r.Header.Get("X-FORWARDED-FOR"); ip != "" {
			ipAddr = ip
		}
		// TODO add more details here
		log.Printf("[SyncAccessLog] %s \"%s %s %s",
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
		case "POST":
			w.Header().Set("Content-Type", "text/plain")
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("[Error] ", err)
				return
			}

			var occurrences []Occurrence
			err = json.Unmarshal(body, &occurrences)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			merge(occurrences, storage)
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})
}

func merge(occurrences []Occurrence, storage Storage) {
	data := make(map[uint64][]uint64)
	for _, v := range occurrences {
		data[v.HashKeyword] = v.Totals
	}
	storage.Merge(data)
}
