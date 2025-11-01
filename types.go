package main

import "time"

type Alert struct {
	DeviceID string    `json:"device_id" bson:"device_id"`
	Estado   string    `json:"estado"`
	PercLOS  float64   `json:"perclos"`
	Blinks   int       `json:"blinks"`
	Yawns    int       `json:"yawns"`
	Ts       time.Time `json:"time"`
}

type HistoryPayload struct {
	Start time.Time `json:"device_id" bson:"device_id"`
	End   time.Time `json:"estado"`
}
