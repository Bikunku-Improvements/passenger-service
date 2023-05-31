package location

import (
	"time"
)

type Event struct {
	BusID     uint      `json:"bus_id"`
	Number    int       `json:"number"`
	Plate     string    `json:"plate"`
	Status    string    `json:"status"`
	Route     string    `json:"route"`
	IsActive  bool      `json:"isActive"`
	Long      float64   `json:"long"`
	Lat       float64   `json:"lat"`
	Speed     float64   `json:"speed"`
	Heading   float64   `json:"heading"`
	CreatedAt time.Time `json:"created_at"`
}
