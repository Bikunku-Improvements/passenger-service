package location

import (
	"encoding/json"
	"time"
)

type Location struct {
	Long      string    `json:"long"`
	Lat       string    `json:"lat"`
	BusID     string    `json:"bus_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (l *Location) UnmarshalBinary(data []byte) error {
	err := json.Unmarshal(data, l)
	if err != nil {
		return err
	}

	return nil
}
