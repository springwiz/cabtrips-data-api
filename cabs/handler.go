package cabs

import (
	"cabtrips-data-api/cache"
	"cabtrips-data-api/model"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type HandlerConfig struct {
	Mysql Repository
	Cache cache.Repository
}

func NewHandlerConfig(Mysql Repository, Cache cache.Repository) HandlerConfig {
	return HandlerConfig{
		Mysql: Mysql,
		Cache: Cache,
	}
}

// GET /cabtrip/{id}
// implements and returns the GET GetCabtripByIDHandler
func GetCabtripByIDHandler(resource HandlerConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		log.Infof("Cabtrip medallion: %s", vars["id"])
		cacheKey := r.FormValue("cache")
		var cabServ CabtripService = newCabtripService(resource.Cache, resource.Mysql)
		trips, err := cabServ.GetCabtripByMedallion(vars["id"], cacheKey)
		if err != nil {
			errRes, _ := json.Marshal(model.NewException("DBERR00001", "Database Error"))
			w.WriteHeader(500)
			err = json.NewEncoder(w).Encode(string(errRes))
			if err != nil {
				return
			}
			return
		}
		if len(trips) == 0 {
			errRes, _ := json.Marshal(model.NewException("DBERR00002", "Medallion not found"))
			w.WriteHeader(404)
			err = json.NewEncoder(w).Encode(string(errRes))
			if err != nil {
				return
			}
			return
		}
		var count model.CabtripCount
		count.Medallion = trips[0].Medallion
		count.TripCount = len(trips)
		countBytes, err := json.Marshal(count)
		if err != nil {
			errRes, _ := json.Marshal(model.NewException("PE00001", "Parse Error"))
			w.WriteHeader(500)
			err = json.NewEncoder(w).Encode(string(errRes))
			if err != nil {
				return
			}
			return
		}
		log.Infof("Response published: %s", string(countBytes))
		w.WriteHeader(200)
		err = json.NewEncoder(w).Encode(string(countBytes))
		if err != nil {
			log.Errorf("Marshal error")
		}
	}
}

//  GET /cabtrip/{id}/date/{pickupdate}
// implements and returns the GET CabtripByPickupdateHandler
func GetCabtripByPickupdateHandler(resource HandlerConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		log.Infof("Cabtrip medallion: %s", vars["id"])
		log.Infof("Cabtrip pickupdate: %s", vars["pickupdate"])
		cacheKey := r.FormValue("cache")
		var cabServ CabtripService = newCabtripService(resource.Cache, resource.Mysql)
		trips, err := cabServ.GetCabtripByMedallionAndPickupdate(vars["id"], vars["pickupdate"], cacheKey)
		if err != nil {
			errRes, _ := json.Marshal(model.NewException("DBERR00001", "Database Error"))
			w.WriteHeader(500)
			err = json.NewEncoder(w).Encode(string(errRes))
			if err != nil {
				return
			}
			return
		}
		if len(trips) == 0 {
			errRes, _ := json.Marshal(model.NewException("DBERR00002", "No trips found for medallion on pickupdate"))
			w.WriteHeader(404)
			err = json.NewEncoder(w).Encode(string(errRes))
			if err != nil {
				return
			}
			return
		}
		var count model.CabtripCount
		count.Medallion = trips[0].Medallion
		count.TripCount = len(trips)
		countBytes, err := json.Marshal(count)
		if err != nil {
			errRes, _ := json.Marshal(model.NewException("PE00001", "Parse Error"))
			w.WriteHeader(500)
			err = json.NewEncoder(w).Encode(string(errRes))
			if err != nil {
				return
			}
			return
		}
		log.Infof("Response published: %s", string(countBytes))
		w.WriteHeader(200)
		err = json.NewEncoder(w).Encode(string(countBytes))
		if err != nil {
			log.Errorf("Marshal error")
		}
	}
}
