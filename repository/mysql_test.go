package repository

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

func TestGetCabtripByMedallion(t *testing.T) {
	tripDb, err := sql.Open("mysql", "root"+":"+"password"+"@tcp("+"localhost"+":"+"3306"+")/"+"cabtrips?parseTime=true")
	if err != nil {
		t.Errorf("error getting db connection %s", err)
	}
	mysql := NewMysql(tripDb)
	cabtrips, err := mysql.GetCabtripByMedallion("9A80FE5419FEA4F44DB8E67F29D84A0F")
	if err != nil {
		t.Errorf("error getting data %s", err)
	}
	assert.True(t, len(cabtrips) > 0)
	assert.True(t, cabtrips[0].Medallion == "9A80FE5419FEA4F44DB8E67F29D84A0F")
	defer tripDb.Close()
}

func TestGetCabtripByMedallionAndPickupdate(t *testing.T) {
	tripDb, err := sql.Open("mysql", "root"+":"+"password"+"@tcp("+"localhost"+":"+"3306"+")/"+"cabtrips?parseTime=true")
	if err != nil {
		t.Errorf("error getting db connection %s", err)
	}
	mysql := NewMysql(tripDb)
	cabtrips, err := mysql.GetCabtripByMedallionAndPickupdate("9A80FE5419FEA4F44DB8E67F29D84A0F", "20131231")
	if err != nil {
		t.Errorf("error getting data %s", err)
	}
	assert.True(t, len(cabtrips) > 0)
	assert.True(t, cabtrips[0].Medallion == "9A80FE5419FEA4F44DB8E67F29D84A0F")
	defer tripDb.Close()
}

func TestRefresh(t *testing.T) {
	tripDb, err := sql.Open("mysql", "root"+":"+"password"+"@tcp("+"localhost"+":"+"3306"+")/"+"cabtrips?parseTime=true")
	if err != nil {
		t.Errorf("error getting db connection %s", err)
	}
	mysql := NewMysql(tripDb)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	err = mysql.Refresh()
	if err != nil {
		log.Fatal(err)
	}
}
