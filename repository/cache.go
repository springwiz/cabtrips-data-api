// Package repository contains all the repository spec and implementations
package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"cabtrips-data-api/model"

	"github.com/allegro/bigcache"
)

// Cache is BigCache based implementation of repository
type Cache struct {
	TripCache *bigcache.BigCache
	TripDb    *sql.DB
}

func NewCache(tripCache *bigcache.BigCache, tripDb *sql.DB) *Cache {
	return &Cache{
		TripCache: tripCache,
		TripDb:    tripDb,
	}
}

func (c *Cache) GetCabtripByMedallion(medallion string) ([]model.Cabtrip, error) {
	var cabtrips []model.Cabtrip
	bytes, err := c.TripCache.Get(medallion)
	if err != nil {
		return nil, fmt.Errorf("no cabtrips found")
	}
	log.Infof("Cache hit with key: %s", medallion)
	json.Unmarshal(bytes, &cabtrips)
	return cabtrips, nil
}

func (c *Cache) GetCabtripByMedallionAndPickupdate(medallion string, pickupDate string) ([]model.Cabtrip, error) {
	var cabtrips []model.Cabtrip
	var ret []model.Cabtrip
	bytes, err := c.TripCache.Get(medallion)
	if err != nil {
		return nil, fmt.Errorf("no cabtrips found")
	}
	log.Infof("Cache hit with key: %s", medallion)
	json.Unmarshal(bytes, &cabtrips)
	for _, cabtrip := range cabtrips {
		strDate := cabtrip.PickupDatetime.Format(`20060102`)
		if strDate == pickupDate {
			ret = append(ret, cabtrip)
		}
	}
	return ret, nil
}

func (c *Cache) Refresh() error {
	rows, err := c.TripDb.Query("select * from cab_trip_data where pickup_datetime >= " +
		"DATE_FORMAT('20131224', '%Y%m%d')")
	if err != nil {
		log.Errorf("error refreshing cache %s", err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var cabtrip model.Cabtrip
		err := rows.Scan(
			&cabtrip.Medallion,
			&cabtrip.HackLicense,
			&cabtrip.VendorID,
			&cabtrip.RateCode,
			&cabtrip.StoreAndFwdFlag,
			&cabtrip.PickupDatetime,
			&cabtrip.DropoffDatetime,
			&cabtrip.PassengerCount,
			&cabtrip.TripTimeInSecs,
			&cabtrip.TripDistance,
			&cabtrip.PickupLongitude,
			&cabtrip.PickupLatitude,
			&cabtrip.DropoffLongitude,
			&cabtrip.DropoffLatitude,
		)
		if err != nil {
			log.Fatal(err)
		}
		if bytes, err := c.TripCache.Get(cabtrip.Medallion); err == nil {
			var cabtrips []model.Cabtrip
			err := json.Unmarshal(bytes, &cabtrips)
			if err != nil {
				log.Fatal(err)
			}
			cabtrips = append(cabtrips, cabtrip)
			bytes, err = json.Marshal(cabtrips)
			if err != nil {
				log.Fatal(err)
			}
			c.TripCache.Set(cabtrip.Medallion, bytes)
		} else {
			var cabtrips []model.Cabtrip
			cabtrips = append(cabtrips, cabtrip)
			bytes, err = json.Marshal(cabtrips)
			if err != nil {
				log.Fatal(err)
			}
			c.TripCache.Set(cabtrip.Medallion, bytes)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
