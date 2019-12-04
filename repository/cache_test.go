package repository

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/allegro/bigcache"
	_ "github.com/go-sql-driver/mysql"

	"github.com/stretchr/testify/assert"
)

func TestGetCabtripByMedallionCache(t *testing.T) {
	t.Parallel()
	tripDb, err := GetMockDB(
		"select \\* from cab_trip_data where medallion = \\?",
		"9A80FE5419FEA4F44DB8E67F29D84A0F")
	if err != nil {
		log.Fatal(err)
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
	mysql := NewMysql(tripDb)
	cabtrips, err := mysql.GetCabtripByMedallion("9A80FE5419FEA4F44DB8E67F29D84A0F")
	if err != nil {
		t.Errorf("error getting data %s", err)
	}
	bytes, err := json.Marshal(cabtrips)
	if err != nil {
		log.Fatal(err)
	}
	err = tripCache.Set("9A80FE5419FEA4F44DB8E67F29D84A0F", bytes)
	if err != nil {
		log.Fatal(err)
	}
	cache := NewCache(tripCache, tripDb)
	cabtrips1, err1 := cache.GetCabtripByMedallion("9A80FE5419FEA4F44DB8E67F29D84A0F")
	if err1 != nil {
		t.Errorf("error getting data %s", err1)
	}
	assert.True(t, len(cabtrips1) > 0)
	assert.True(t, cabtrips1[0].Medallion == "9A80FE5419FEA4F44DB8E67F29D84A0F")
	defer tripDb.Close()
}

func TestGetCabtripByMedallionAndPickupdateCache(t *testing.T) {
	t.Parallel()
	tripDb, err := GetMockDB(
		"select \\* from cab_trip_data where medallion = \\? and DATE_FORMAT\\(pickup_datetime, '%Y%m%d'\\) = \\?",
		"9A80FE5419FEA4F44DB8E67F29D84A0F",
		"20131231")
	if err != nil {
		log.Fatal(err)
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
	mysql := NewMysql(tripDb)
	cabtrips, err := mysql.GetCabtripByMedallionAndPickupdate("9A80FE5419FEA4F44DB8E67F29D84A0F", "20131231")
	if err != nil {
		t.Errorf("error getting data %s", err)
	}
	bytes, err := json.Marshal(cabtrips)
	if err != nil {
		log.Fatal(err)
	}
	err = tripCache.Set("9A80FE5419FEA4F44DB8E67F29D84A0F", bytes)
	if err != nil {
		log.Fatal(err)
	}
	cache := NewCache(tripCache, tripDb)
	cabtrips1, err1 := cache.GetCabtripByMedallionAndPickupdate("9A80FE5419FEA4F44DB8E67F29D84A0F", "20131231")
	if err1 != nil {
		t.Errorf("error getting data %s", err1)
	}
	assert.True(t, len(cabtrips1) > 0)
	assert.True(t, cabtrips1[0].Medallion == "9A80FE5419FEA4F44DB8E67F29D84A0F")
	defer tripDb.Close()
}

func TestRefreshCache(t *testing.T) {
	t.Parallel()
	tripDb, err := GetMockDB(
		"select \\* from cab_trip_data where pickup_datetime >= DATE_FORMAT\\('20131224', '%Y%m%d'\\)")
	if err != nil {
		log.Fatal(err)
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
	cache := NewCache(tripCache, tripDb)
	err = cache.Refresh()
	if err != nil {
		log.Fatal(err)
	}
}
