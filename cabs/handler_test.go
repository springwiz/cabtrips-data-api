package cabs

import (
	"cabtrips-data-api/repository"
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"github.com/allegro/bigcache"
)

// GET /cabtrip/{id}
// Tests the GET GetCabtripByIDHandler
func TestGetCabtripByIDHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/cabtrip", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "9A80FE5419FEA4F44DB8E67F29D84A0F"})
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	res, err := getTestResource()
	if err != nil {
		t.Fatal(err)
	}
	handler := http.HandlerFunc(GetCabtripByIDHandler(*res))
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := regexp.MustCompile(`{\"medallion\":\"9A80FE5419FEA4F44DB8E67F29D84A0F\",\"tripCount\":2}`)
	if matches := expected.MatchString(rr.Body.String()); matches {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

//  GET /cabtrip/{id}/date/{pickupdate}
// Tests the GET GetCabtripByPickupdateHandler
func TestGetCabtripByPickupdateHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/cabtrip", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id":         "9A80FE5419FEA4F44DB8E67F29D84A0F",
		"pickupdate": "20131231"})
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	res, err := getTestResource()
	if err != nil {
		t.Fatal(err)
	}
	handler := http.HandlerFunc(GetCabtripByPickupdateHandler(*res))
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := regexp.MustCompile(`{\"medallion\":\"9A80FE5419FEA4F44DB8E67F29D84A0F\",\"tripCount\":1}`)
	if matches := expected.MatchString(rr.Body.String()); matches {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetCabtripByIDCacheHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/cabtrip", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id":    "9A80FE5419FEA4F44DB8E67F29D84A0F",
		"cache": "true",
	})
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	res, err := getTestResource()
	if err != nil {
		t.Fatal(err)
	}
	handler := http.HandlerFunc(GetCabtripByIDHandler(*res))
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := regexp.MustCompile(`{\"medallion\":\"9A80FE5419FEA4F44DB8E67F29D84A0F\",\"tripCount\":2}`)
	if matches := expected.MatchString(rr.Body.String()); matches {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

//  GET /cabtrip/{id}/date/{pickupdate}
// Tests the GET GetCabtripByPickupdateHandler
func TestGetCabtripByPickupdateCacheHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/cabtrip", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id":         "9A80FE5419FEA4F44DB8E67F29D84A0F",
		"pickupdate": "20131231",
		"cache":      "true",
	})
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	res, err := getTestResource()
	if err != nil {
		t.Fatal(err)
	}
	handler := http.HandlerFunc(GetCabtripByPickupdateHandler(*res))
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := regexp.MustCompile(`{\"medallion\":\"9A80FE5419FEA4F44DB8E67F29D84A0F\",\"tripCount\":1}`)
	if matches := expected.MatchString(rr.Body.String()); matches {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func getTestResource() (*HandlerConfig, error) {
	tripDb, err := sql.Open("mysql", "root"+":"+"password"+"@tcp("+"localhost"+":"+"3306"+")/"+"cabtrips?parseTime=true")
	if err != nil {
		return nil, err
	}
	mysql := repository.NewMysql(tripDb)
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
		Mysql: mysql,
		Cache: repository.NewCache(tripCache, tripDb),
	}, nil
}
