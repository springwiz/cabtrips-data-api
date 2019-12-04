// Package repository contains all the repository spec and implementations
package repository

import (
	"cabtrips-data-api/model"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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

func GetMockDB(mockquery string, mockargs ...driver.Value) (*sql.DB, error) {
	// create sqlmock database connection and a mock to manage expectations.
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, fmt.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	var pTime, dTime time.Time
	pTime, err = time.Parse("2006-01-02 15:04:05", "2013-12-31 07:39:00")
	if err != nil {
		return nil, fmt.Errorf("time parse error %s", err)
	}
	dTime, err = time.Parse("2006-01-02 15:04:05", "2013-12-31 07:46:00")
	if err != nil {
		return nil, fmt.Errorf("time parse error %s", err)
	}
	rows := sqlmock.NewRows([]string{
		"medallion",
		"hack_license",
		"vendor_id",
		"rate_code",
		"store_and_fwd_flag",
		"pickup_datetime",
		"dropoff_datetime",
		"passenger_count",
		"trip_time_in_secs",
		"trip_distance",
		"pickup_longitude",
		"pickup_latitude",
		"dropoff_longitude",
		"dropoff_latitude"}).AddRow(
		"9A80FE5419FEA4F44DB8E67F29D84A0F",
		"-73.972794",
		"VTS",
		"1",
		"-73.995262",
		pTime,
		dTime,
		"5", "420", "2.29",
		nil, nil, nil, nil)
	mockObj := mock.ExpectQuery(mockquery)
	mockObj.WithArgs(mockargs...)
	mockObj.WillReturnRows(rows)
	return db, nil
}
