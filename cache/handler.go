package cache

import (
	"cabtrips-data-api/model"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type HandlerConfig struct {
	Cache Repository
}

func NewHandlerConfig(Cache Repository) HandlerConfig {
	return HandlerConfig{
		Cache: Cache,
	}
}

//  GET /cabtrip/refresh_cache
// implements and returns the GET GetRefreshCacheHandler
func GetRefreshCacheHandler(resource HandlerConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := resource.Cache.Refresh()
		if err != nil {
			errRes, _ := json.Marshal(model.NewException("DBERR00001", "Database Error"))
			w.WriteHeader(500)
			err = json.NewEncoder(w).Encode(string(errRes))
			if err != nil {
				return
			}
			return
		}
		log.Info("Cache refreshed")
		w.WriteHeader(200)
	}
}
