// Package repository contains all the repository spec and implementations
package repository

import (
	"cabtrips-data-api/model"
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// Mysql is Mysql based implementation of repository
type Mysql struct {
	TripDb *sql.DB
}

func NewMysql(tripDb *sql.DB) *Mysql {
	return &Mysql{
		TripDb: tripDb,
	}
}

func (m *Mysql) GetCabtripByMedallion(medallion string) ([]model.Cabtrip, error) {
	rows, err := m.TripDb.Query("select * from cab_trip_data where medallion = ?", medallion)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return m.copyRows(rows)
}

func (m *Mysql) GetCabtripByMedallionAndPickupdate(medallion string, pickupDate string) ([]model.Cabtrip, error) {
	rows, err := m.TripDb.Query("select * from cab_trip_data where medallion = ? and "+
		"DATE_FORMAT(pickup_datetime, '%Y%m%d') = ?", medallion, pickupDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return m.copyRows(rows)
}

func (m *Mysql) copyRows(rows *sql.Rows) ([]model.Cabtrip, error) {
	var cabtrips []model.Cabtrip
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
		cabtrips = append(cabtrips, cabtrip)
	}
	return cabtrips, nil
}

func (m *Mysql) Refresh() error {
	panic("operation usupported")
}
