package models

import "time"

type Trade struct {
	Time   time.Time
	Signal string
	Price  float64
}
