package cache

import (
	"cabtrips-data-api/repository"
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/allegro/bigcache"
)

//  GET /cabtrip/refresh_cache
// Tests the GET GetRefreshCacheHandler
func TestGetRefreshCacheHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/cabtrip/refresh_cache", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	res, err := getTestResource()
	if err != nil {
		t.Fatal(err)
	}
	handler := http.HandlerFunc(GetRefreshCacheHandler(*res))
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func getTestResource() (*HandlerConfig, error) {
	tripDb, err := sql.Open("mysql", "root"+":"+"password"+"@tcp("+"localhost"+":"+"3306"+")/"+"cabtrips?parseTime=true")
	if err != nil {
		return nil, err
	}
	config := bigcache.Config{
		Shards:             1024,
		LifeWindow:         10,
		MaxEntriesInWindow: 600000,
		MaxEntrySize:       500,
		Verbose:            true,
		HardMaxCacheSize:   0,
	}
	tripCache, err := bigcache.NewBigCache(config)
	if err != nil {
		log.Fatal(err)
	}
	return &HandlerConfig{
		Cache: repository.NewCache(tripCache, tripDb),
	}, nil
}
