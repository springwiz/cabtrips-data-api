package main

import (
	"cabtrips-data-api/cabs"
	"cabtrips-data-api/cache"
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

type config struct {
	Mysql      cabs.Repository
	Cache      cache.Repository
	Host, Port string
}

func initialize() config {
	var serviceConfig config
	// read the config yml
	viper.SetConfigName("server")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Warnf("Config file not found...")
		log.Warnf("Using Defaults")
		serviceConfig.Host = "localhost"
		serviceConfig.Port = "8080"
	} else {
		serviceConfig.Host = viper.GetString("server.host")
		serviceConfig.Port = viper.GetString("server.port")
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
	serviceConfig.Cache = repository.NewCache(tripCache, tripDb)
	serviceConfig.Mysql = repository.NewMysql(tripDb)
	err = serviceConfig.Cache.Refresh()
	if err != nil {
		log.Fatalf("cache refresh failed %s", err)
	}
	return serviceConfig
}

// serves as the starting point of the application
// initializes the mux router and maps the routes to functions
func main() {
	r := mux.NewRouter()
	// acquire and initialize the resoources
	serviceConfig := initialize()
	cabtripHandlerConfig := cabs.NewHandlerConfig(serviceConfig.Mysql, serviceConfig.Cache)
	cacheHandlerConfig := cache.NewHandlerConfig(serviceConfig.Cache)
	r.HandleFunc("/cabtrip/{id}",
		cabs.GetCabtripByIDHandler(cabtripHandlerConfig)).Methods("GET")
	r.HandleFunc("/cabtrip/{id}",
		cabs.GetCabtripByIDHandler(cabtripHandlerConfig)).Queries("cache", "").Methods("GET")
	r.HandleFunc("/cabtrip/{id}/date/{pickupdate}",
		cabs.GetCabtripByPickupdateHandler(cabtripHandlerConfig)).Methods("GET")
	r.HandleFunc("/cabtrip/{id}/date/{pickupdate}",
		cabs.GetCabtripByPickupdateHandler(cabtripHandlerConfig)).Queries("cache", "").Methods("GET")
	r.HandleFunc("/cache/refresh_cache", cache.GetRefreshCacheHandler(cacheHandlerConfig)).Methods("GET")
	err := http.ListenAndServe(serviceConfig.Host+":"+serviceConfig.Port, r)
	if err != nil {
		log.Fatalf("starting server failed %s", err)
	}
}
