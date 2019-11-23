// Package handler provides with handler functions for handling the various HTTP Requests
package handler

import (
	"cabtrips-data-api/model"
	"cabtrips-data-api/repository"
	"cabtrips-data-api/service"

	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Resource struct {
	Mysql      *repository.Mysql
	Cache      *repository.Cache
	Host, Port string
}

// GET /cabtrip/{id}
// implements and returns the GET GetCabtripByIDHandler
func GetCabtripByIDHandler(resource Resource) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		log.Infof("Cabtrip medallion: %s", vars["id"])
		cacheKey := r.FormValue("cache")
		var cabServ service.Cabtrip = service.NewCabtrip(resource.Cache, resource.Mysql)
		trips, err := cabServ.GetCabtripByMedallion(vars["id"], cacheKey)
		if err != nil {
			errRes, _ := json.Marshal(model.NewException("DBERR00001", "Database Error"))
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(string(errRes))
			return
		}
		if len(trips) == 0 {
			errRes, _ := json.Marshal(model.NewException("DBERR00002", "Medallion not found"))
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(string(errRes))
			return
		}
		var count model.CabtripCount
		count.Medallion = trips[0].Medallion
		count.TripCount = len(trips)
		countBytes, err := json.Marshal(count)
		if err != nil {
			errRes, _ := json.Marshal(model.NewException("PE00001", "Parse Error"))
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(string(errRes))
			return
		}
		log.Infof("Response published: %s", string(countBytes))
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(string(countBytes))
	}
}

//  GET /cabtrip/{id}/date/{pickupdate}
// implements and returns the GET CabtripByPickupdateHandler
func GetCabtripByPickupdateHandler(resource Resource) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		log.Infof("Cabtrip medallion: %s", vars["id"])
		log.Infof("Cabtrip pickupdate: %s", vars["pickupdate"])
		cacheKey := r.FormValue("cache")
		var cabServ service.Cabtrip = service.NewCabtrip(resource.Cache, resource.Mysql)
		trips, err := cabServ.GetCabtripByMedallionAndPickupdate(vars["id"], vars["pickupdate"], cacheKey)
		if err != nil {
			errRes, _ := json.Marshal(model.NewException("DBERR00001", "Database Error"))
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(string(errRes))
			return
		}
		if len(trips) == 0 {
			errRes, _ := json.Marshal(model.NewException("DBERR00002", "No trips found for medallion on pickupdate"))
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(string(errRes))
			return
		}
		var count model.CabtripCount
		count.Medallion = trips[0].Medallion
		count.TripCount = len(trips)
		countBytes, err := json.Marshal(count)
		if err != nil {
			errRes, _ := json.Marshal(model.NewException("PE00001", "Parse Error"))
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(string(errRes))
			return
		}
		log.Infof("Response published: %s", string(countBytes))
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(string(countBytes))
	}
}

//  GET /cabtrip/refresh_cache
// implements and returns the GET GetRefreshCacheHandler
func GetRefreshCacheHandler(resource Resource) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := resource.Cache.Refresh()
		if err != nil {
			errRes, _ := json.Marshal(model.NewException("DBERR00001", "Database Error"))
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(string(errRes))
			return
		}
		log.Info("Cache refreshed")
		w.WriteHeader(200)
	}
}
