package main

import (
	"cabtrips-data-api/handler"
	"cabtrips-data-api/repository"
	"database/sql"
	"net/http"
	"time"

	"github.com/allegro/bigcache"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var resource handler.Resource

func init() {
	// read the config yml
	viper.SetConfigName("server")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Warnf("Config file not found...")
		log.Warnf("Using Defaults")
		resource.Host = "localhost"
		resource.Port = "8080"
	} else {
		resource.Host = viper.GetString("server.host")
		resource.Port = viper.GetString("server.port")
	}
	config := bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: viper.GetInt("bigcache.shards"),
		// time after which entry can be evicted
		LifeWindow: viper.GetDuration("bigcache.lifeWindow") * time.Minute,
		// rps * lifeWindow, used only in initial memory allocation
		MaxEntriesInWindow: viper.GetInt("bigcache.maxEntriesInWindow"),
		// max entry size in bytes, used only in initial memory allocation
		MaxEntrySize: viper.GetInt("bigcache.maxEntrySize"),
		// prints information about additional memory allocation
		Verbose: viper.GetBool("bigcache.verbose"),
		// cache will not allocate more memory than this limit, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		HardMaxCacheSize: viper.GetInt("bigcache.hardMaxCacheSize"),
	}
	tripCache, err := bigcache.NewBigCache(config)
	if err != nil {
		log.Fatal(err)
	}
	sqlString := viper.GetString("mysql.user") + ":" + viper.GetString("mysql.password") +
		"@tcp(" + viper.GetString("mysql.host") + ":" + viper.GetString("mysql.port") +
		")/" + viper.GetString("mysql.schema") + "?parseTime=true"
	tripDb, err := sql.Open("mysql", sqlString)
	if err != nil {
		log.Fatal(err)
	}
	resource.Cache = repository.NewCache(tripCache, tripDb)
	resource.Mysql = repository.NewMysql(tripDb)
	resource.Cache.Refresh()
}

// serves as the starting point of the application
// initializes the mux router and maps the routes to functions
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/cabtrip/{id}", handler.GetCabtripByIDHandler(resource)).Methods("GET")
	r.HandleFunc("/cabtrip/{id}", handler.GetCabtripByIDHandler(resource)).Queries("cache", "").Methods("GET")
	r.HandleFunc("/cabtrip/{id}/date/{pickupdate}", handler.GetCabtripByPickupdateHandler(resource)).Methods("GET")
	r.HandleFunc("/cabtrip/{id}/date/{pickupdate}", handler.GetCabtripByPickupdateHandler(resource)).
		Queries("cache", "").Methods("GET")
	r.HandleFunc("/cache/refresh_cache", handler.GetRefreshCacheHandler(resource)).Methods("GET")
	http.ListenAndServe(resource.Host+":"+resource.Port, r)
}
