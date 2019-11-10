// Package model contains all the entities
package model

// The CabtripCount entity
type CabtripCount struct {
	Medallion string `json:"medallion"`

	TripCount int `json:"tripCount"`
}
